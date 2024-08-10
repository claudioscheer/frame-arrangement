package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// Wall dimensions
type Wall struct {
	width  int
	height int
}

// Frame dimensions
type Frame struct {
	width  int
	height int
	x      int // x position
	y      int // y position
}

// Checks if the frame overlaps with any existing frames
func hasCollision(newFrame Frame, frames []Frame, margin int) bool {
	for _, frame := range frames {
		if newFrame.x < frame.x+frame.width+margin &&
			newFrame.x+newFrame.width+margin > frame.x &&
			newFrame.y < frame.y+frame.height+margin &&
			newFrame.y+newFrame.height+margin > frame.y {
			return true
		}
	}
	return false
}

// Tries to place frames randomly near previous frames until it fills 85% of the wall
func placeFrames(wall Wall, frameSizes []Frame, marginRange [2]int, rng *rand.Rand) []Frame {
	var placedFrames []Frame
	requiredArea := int(float64(wall.width*wall.height) * 0.54)
	totalArea := 0

	// Place the first frame randomly on the wall
	firstFrame := frameSizes[0]
	firstFrame.x, firstFrame.y = rng.Intn(wall.width-firstFrame.width), rng.Intn(wall.height-firstFrame.height)
	placedFrames = append(placedFrames, firstFrame)
	totalArea += firstFrame.width * firstFrame.height

	// Place subsequent frames near previous frames with some random offset
	for totalArea < requiredArea {
		fmt.Println("Total Area in Percentage: ", (totalArea*100)/(wall.width*wall.height))
		rng.Shuffle(len(frameSizes), func(i, j int) {
			frameSizes[i], frameSizes[j] = frameSizes[j], frameSizes[i]
		})

		for _, frame := range frameSizes {
			if totalArea >= requiredArea {
				break
			}

			margin := rng.Intn(marginRange[1]-marginRange[0]) + marginRange[0]
			placed := false

			// Try to place the frame near each existing frame with more candidate positions
			for _, prevFrame := range placedFrames {
				// Generate 12 potential positions around the existing frame
				candidates := []Frame{
					{width: frame.width, height: frame.height, x: prevFrame.x - frame.width - margin, y: prevFrame.y},                                 // Left
					{width: frame.width, height: frame.height, x: prevFrame.x + prevFrame.width + margin, y: prevFrame.y},                             // Right
					{width: frame.width, height: frame.height, x: prevFrame.x, y: prevFrame.y - frame.height - margin},                                // Above
					{width: frame.width, height: frame.height, x: prevFrame.x, y: prevFrame.y + prevFrame.height + margin},                            // Below
					{width: frame.width, height: frame.height, x: prevFrame.x - frame.width - margin, y: prevFrame.y - margin},                        // Top-left
					{width: frame.width, height: frame.height, x: prevFrame.x + prevFrame.width + margin, y: prevFrame.y - margin},                    // Top-right
					{width: frame.width, height: frame.height, x: prevFrame.x - frame.width - margin, y: prevFrame.y + prevFrame.height + margin},     // Bottom-left
					{width: frame.width, height: frame.height, x: prevFrame.x + prevFrame.width + margin, y: prevFrame.y + prevFrame.height + margin}, // Bottom-right
					{width: frame.width, height: frame.height, x: prevFrame.x - frame.width, y: prevFrame.y - frame.height},                           // Top-left diagonal
					{width: frame.width, height: frame.height, x: prevFrame.x + prevFrame.width, y: prevFrame.y - frame.height},                       // Top-right diagonal
					{width: frame.width, height: frame.height, x: prevFrame.x - frame.width, y: prevFrame.y + prevFrame.height},                       // Bottom-left diagonal
					{width: frame.width, height: frame.height, x: prevFrame.x + prevFrame.width, y: prevFrame.y + prevFrame.height},                   // Bottom-right diagonal
				}

				// Shuffle the candidates and try placing them
				rng.Shuffle(len(candidates), func(i, j int) {
					candidates[i], candidates[j] = candidates[j], candidates[i]
				})

				for _, candidate := range candidates {
					// Check bounds and collisions
					if candidate.x >= 0 && candidate.y >= 0 &&
						candidate.x+candidate.width <= wall.width &&
						candidate.y+candidate.height <= wall.height &&
						!hasCollision(candidate, placedFrames, margin) {
						placedFrames = append(placedFrames, candidate)
						totalArea += candidate.width * candidate.height
						placed = true
						break
					}
				}

				if placed {
					break
				}
			}

			// If no valid position was found, try another frame
			if !placed {
				continue
			}
		}
	}

	return placedFrames
}

// Draw the frames on the wall image and save as a PNG
func visualize(wall Wall, frames []Frame) {
	// Create a white background
	img := image.NewRGBA(image.Rect(0, 0, wall.width, wall.height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Draw each frame as a colored rectangle
	for i, frame := range frames {
		frameColor := color.RGBA{uint8(100 + i*20), uint8(50 + i*15), uint8(150 + i*10), 255}
		for x := frame.x; x < frame.x+frame.width; x++ {
			for y := frame.y; y < frame.y+frame.height; y++ {
				img.Set(x, y, frameColor)
			}
		}
	}

	// Save the image to a file
	file, err := os.Create("wall_visualization.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}

	fmt.Println("Visualization saved as wall_visualization.png")
}

func main() {
	// Create a new random number generator with a time-based seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	wall := Wall{width: 230, height: 140}
	frameSizes := []Frame{
		{width: 10, height: 15},
		{width: 15, height: 10},
		{width: 13, height: 18},
		{width: 18, height: 13},
		{width: 16, height: 9},
		{width: 9, height: 9},
	}

	marginRange := [2]int{2, 5} // Reduced margin range for tighter packing

	placedFrames := placeFrames(wall, frameSizes, marginRange, rng)

	visualize(wall, placedFrames)
}
