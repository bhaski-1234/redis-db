package command

func HandlePing(args []string) (interface{}, error) {
	if len(args) == 1 {
		return "PONG", nil
	}
	return args[1], nil
}
