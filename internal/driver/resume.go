package driver

import "github.com/makasim/flowstate"

var ResumeDoer flowstate.DoerFunc = func(cmd0 flowstate.Command) error {
	if err := flowstate.DefaultResumeDoer.Do(cmd0); err != nil {
		return err
	}

	cmd, ok := cmd0.(*flowstate.ResumeCommand)
	if !ok {
		return flowstate.ErrCommandNotSupported
	}

	cmd.StateCtx.Current.SetLabel(`flowstate.state_id`, string(cmd.StateCtx.Current.ID))
	cmd.StateCtx.Current.SetLabel(`flowstate.flow_id`, string(cmd.StateCtx.Current.Transition.ToID))
	return nil
}