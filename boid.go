package main

import (
	"math"
	"math/rand"
	"time"
)

// Boid is kind of like "bot + bird"
type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) clcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0

	rWLock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBid := boidMap[int(i)][int(j)]; otherBid != -1 && otherBid != b.id {
				if dist := boids[otherBid].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBid].velocity)
					avgPosition = avgPosition.Add(boids[otherBid].position)
					separation = separation.Add(b.position.Substract(boids[otherBid].position).DivideV(dist))
				}
			}
		}
	}
	rWLock.RUnlock()

	acc := Vector2D{
		b.borderBounce(b.position.x, screenWidth),
		b.borderBounce(b.position.y, screenHeight),
	}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivideV(count), avgVelocity.DivideV(count)
		accAlignment := avgVelocity.Substract(b.velocity).MultiplyV(adjRate)
		accCohesion := avgPosition.Substract(b.position).MultiplyV(adjRate)
		accSeparation := separation.MultiplyV(adjRate)
		acc = acc.Add(accAlignment).Add(accCohesion).Add(accSeparation)
	}
	return acc
}

func (b *Boid) borderBounce(pos, maxPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxPos-viewRadius {
		return 1 / (pos - maxPos)
	}
	return 0
}

func (b *Boid) moveOne() {
	acc := b.clcAcceleration()

	rWLock.Lock()
	b.velocity = b.velocity.Add(acc).Limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	// next := b.position.Add(b.velocity)
	// if next.x >= screenWidth || next.x < 0 {
	// 	b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	// }
	// if next.y >= screenHeight || next.y < 0 {
	// 	b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	// }
	rWLock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}
