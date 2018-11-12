/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package runc

import (
	"context"
	"os"
	"os/exec"
	"syscall"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

func (r *Runc) command(context context.Context, args ...string) *exec.Cmd {

	// Create an register a OpenCensus
	// Stackdriver Trace exporter.
	exporter, _ := stackdriver.NewExporter(stackdriver.Options{})

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	_, runcStartSpan := trace.StartSpan(context, "RuncStartSpan.Start")
	runcStartSpan.AddAttributes(trace.StringAttribute("command", r.Command))

	command := r.Command
	if command == "" {
		command = DefaultCommand
	}
	cmd := exec.CommandContext(context, command, append(r.args(), args...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: r.Setpgid,
	}
	cmd.Env = os.Environ()
	if r.PdeathSignal != 0 {
		cmd.SysProcAttr.Pdeathsig = r.PdeathSignal
	}

	runcStartSpan.End()

	return cmd
}
