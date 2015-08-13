package commands

// CatchUp marks all episodes downloaded = true
func (env *Env) CatchUp(id int) {
	env.db.CatchUp(id)
}
