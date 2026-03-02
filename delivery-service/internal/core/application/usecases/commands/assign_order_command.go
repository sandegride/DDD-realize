package commands

type AssignOrdersCommand struct {
	valid bool
}

func NewAssignOrdersCommand() (AssignOrdersCommand, error) {

	return AssignOrdersCommand{valid: true}, nil
}

func (cmd AssignOrdersCommand) IsValid() bool {
	return cmd.valid
}
