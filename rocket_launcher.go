package server

//LauncherRange ...
const LauncherRange = 250

//Launcher ...
type Launcher struct {
	pos         *Vector
	reloadTime  int
	reloadTimer int
	player      *Player
}

//NewLauncher ...
func NewLauncher(pos *Vector, reloadTime int, player *Player) *Launcher {
	return &Launcher{
		pos:        pos,
		reloadTime: reloadTime,
		player:     player,
	}
}

//Tick ...
func (l *Launcher) Tick() {
	if l.reloadTime > l.reloadTimer {
		l.reloadTimer++
		return
	}

	for _, otherPlayer := range OtherPlayers(l.player) {
		rocketTargets := otherPlayer.RocketTargets()

		for _, target := range rocketTargets {
			if target.Position().Dist(l.pos) < LauncherRange {
				l.player.AddRocket(l.pos, target)
				l.reloadTimer = 0
				return
			}
		}
	}
}
