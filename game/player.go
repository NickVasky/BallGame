package game

type PlayerController struct {
	movementSpeed float64
	movementDecay float64
	controlVector Vec2
}

func NewPlayerController(speed, decay float64) PlayerController {
	return PlayerController{
		movementSpeed: speed,
		movementDecay: decay,
	}
}
