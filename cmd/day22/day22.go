package main

import (
    "fmt"
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Volume struct {
    xMin, xMax, yMin, yMax, zMin, zMax int
}

func (v Volume) String() string {
    return fmt.Sprintf("%d..%d,%d..%d,%d..%d", v.xMin, v.xMax, v.yMin, v.yMax, v.zMin, v.zMax)
}

func (v Volume) Size() int {
    return (v.xMax - v.xMin + 1) * (v.yMax - v.yMin + 1) * (v.zMax - v.zMin + 1)
}

func (v Volume) Intersects(ov Volume) bool {
    return ! (v.xMax < ov.xMin || v.xMin > ov.xMax || v.yMax < ov.yMin || v.yMin > ov.yMax || v.zMax < ov.zMin || v.zMin > ov.zMax)
}

func (v Volume) Contains(ov Volume) bool {
    return v.xMin <= ov.xMin && ov.xMax <= v.xMax && v.yMin <= ov.yMin && ov.yMax <= v.yMax && v.zMin <= ov.zMin && ov.zMax <= v.zMax
}

func (v Volume) ShearOff(ov Volume, instruct string) Cubes {
    // v is a cube possibly intersecting with cube ov; v's state will "dominate" and the subvolumes of ov not in v will remain in
    // their state. Result should be 1-6 non-interesecting volumes
    if v.xMax < ov.xMin || v.xMin > ov.xMax || v.yMax < ov.yMin || v.yMin > ov.yMax || v.zMax < ov.zMin || v.zMin > ov.zMax {
        // no intersection
        return Cubes{ov: instruct}
    }
    if v.xMax >= ov.xMax && v.xMin <= ov.xMin && v.yMax >= ov.yMax && v.yMin <= ov.yMin && v.zMax >= ov.zMax && v.zMin <= ov.zMin {
        // v encloses ov
        return make(Cubes)
    }
    cubes := make(Cubes)
    yMinOff, yMaxOff, zMinOff, zMaxOff := 0, 0, 0, 0
    if ov.yMin < v.yMin { // bottom xz-plane is uninterrupted
        // shear off ov.yMin up to v.yMin - 1
        cubes[Volume{ov.xMin, ov.xMax, ov.yMin, v.yMin - 1, ov.zMin, ov.zMax}] = instruct
        yMinOff += (v.yMin - ov.yMin)
    }
    if ov.yMax > v.yMax { // top xz-plane is uninterrupted
        // shear off v.yMax + 1 up to ov.yMax
        cubes[Volume{ov.xMin, ov.xMax, v.yMax + 1, ov.yMax, ov.zMin, ov.zMax}] = instruct
        yMaxOff += (ov.yMax - v.yMax)
    }
    if ov.zMin < v.zMin { // rear xy-plane is uninterrupted
        // shear off ov.zMin forward to v.zMin - 1 
        cubes[Volume{ov.xMin, ov.xMax, ov.yMin + yMinOff, ov.yMax - yMaxOff, ov.zMin, v.zMin - 1}] = instruct
        zMinOff += (v.zMin - ov.zMin)
    }
    if ov.zMax > v.zMax { // front xy-plane is uninterrupted
        // shear off v.zMax + 1 forward to ov.zMax
        cubes[Volume{ov.xMin, ov.xMax, ov.yMin + yMinOff, ov.yMax - yMaxOff, v.zMax + 1, ov.zMax}] = instruct
        zMaxOff += (ov.zMax - v.zMax)
    }
    if ov.xMin < v.xMin { // left yz-plane is uninterrupted
        // shear off ov.xMin forward to v.xMin - 1 
        cubes[Volume{ov.xMin, v.xMin - 1, ov.yMin + yMinOff, ov.yMax - yMaxOff, ov.zMin + zMinOff, ov.zMax - zMaxOff}] = instruct
    }
    if ov.xMax > v.xMax { // right yz-plane is uninterrupted
        // shear off v.xMax + 1 forward to ov.xMax
        cubes[Volume{v.xMax + 1, ov.xMax, ov.yMin + yMinOff, ov.yMax - yMaxOff, ov.zMin + zMinOff, ov.zMax - zMaxOff}] = instruct
    }
    return cubes
}

type Cubes map[Volume]string

func cubesFromInput(lines []string) ([]Volume, []string) {
    cubes := make([]Volume, 0)
    instructs := make([]string, 0)
    for _, line := range lines {
        reg := regexp.MustCompile(`(\w+)\sx=([0-9-]+)..([0-9-]+),y=([0-9-]+)..([0-9-]+),z=([0-9-]+)..([0-9-]+)`)
        if result := reg.FindStringSubmatch(line); result != nil {
            instruction := result[1]
            xMinStr := result[2]
            xMin, _ := strconv.Atoi(xMinStr)
            xMaxStr := result[3]
            xMax, _ := strconv.Atoi(xMaxStr)
            yMinStr := result[4]
            yMin, _ := strconv.Atoi(yMinStr)
            yMaxStr := result[5]
            yMax, _ := strconv.Atoi(yMaxStr)
            zMinStr := result[6]
            zMin, _ := strconv.Atoi(zMinStr)
            zMaxStr := result[7]
            zMax, _ := strconv.Atoi(zMaxStr)
            cube := Volume{xMin, xMax, yMin, yMax, zMin, zMax}
            cubes = append(cubes, cube)
            instructs = append(instructs, instruction)
        }
    }
    return cubes, instructs
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    initialMin, initialMax := -50, 50
    initial := Volume{initialMin, initialMax, initialMin, initialMax, initialMin, initialMax}
    cubes, instructs := cubesFromInput(lines)
    vList := make(Cubes)
    for i, cube := range cubes {
        instruct := instructs[i]
        if initial.Contains(cube) {
            // do the thing
            if len(vList) == 0 {
                vList[cube] = instruct
                continue
            }
            newVList := make(Cubes)
            for vol, state := range vList {
                if cube.Intersects(vol) {
                    shears := cube.ShearOff(vol, state)
                    shearSum := 0
                    for shear, shearState := range shears {
                        shearSum += shear.Size()
                        newVList[shear] = shearState
                    }
                    if shearSum >= vol.Size() {
                        panic(fmt.Sprintf("Sum of shears (%d) >= size of original volume (%d)!", shearSum, vol.Size()))
                    }
                } else {
                    newVList[vol] = state
                }
            }
            newVList[cube] = instruct
            vList = newVList
            loopSum := 0
            for vol, state := range vList {
                if state == "on" {
                    loopSum += vol.Size()
                }
            }
        }
    }
    sum := 0
    for cube, instruct := range vList {
        if initial.Contains(cube) {
            if instruct == "on" {
                sum += cube.Size()
            }
        }
    }
    return sum
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    cubes, instructs := cubesFromInput(lines)
    vList := make(Cubes)
    for i, cube := range cubes {
        instruct := instructs[i]
        // do the thing
        if len(vList) == 0 {
            vList[cube] = instruct
            continue
        }
        newVList := make(Cubes)
        for vol, state := range vList {
            if cube.Intersects(vol) {
                shears := cube.ShearOff(vol, state)
                shearSum := 0
                for shear, shearState := range shears {
                    shearSum += shear.Size()
                    newVList[shear] = shearState
                }
                if shearSum >= vol.Size() {
                    panic(fmt.Sprintf("Sum of shears (%d) >= size of original volume (%d)!", shearSum, vol.Size()))
                }
            } else {
                newVList[vol] = state
            }
        }
        newVList[cube] = instruct
        vList = newVList
        loopSum := 0
        for vol, state := range vList {
            if state == "on" {
                loopSum += vol.Size()
            }
        }
    }
    sum := 0
    for cube, instruct := range vList {
        if instruct == "on" {
            sum += cube.Size()
        }
    }
    return sum
}

