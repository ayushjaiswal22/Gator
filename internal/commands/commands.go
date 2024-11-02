package commands

import(
    "errors"
)

func (c *Commands) Register(name string, handlerFunc func(*State, Command) error) {
    c.CmdMap[name] = handlerFunc
}

func (c *Commands) Run(name string, s *State, cmd Command) error {
    handler, ok := c.CmdMap[name]
    if ok {
        err := handler(s, cmd)
        if err!=nil {
            return err
        } else{
            return nil
        }
    } else {
        return errors.New("Invalid command")
    }
}
