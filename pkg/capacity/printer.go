// Copyright 2019 Rob Scott
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capacity

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

func printList(cm *clusterMetric, showPods bool, showUtil bool) {
	names := make([]string, len(cm.nodeMetrics))

	i := 0
	for name := range cm.nodeMetrics {
		names[i] = name
		i++
	}
	sort.Strings(names)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)

	printHeaders(w, cm, showPods, showUtil)

	for _, name := range names {
		printNode(w, name, cm.nodeMetrics[name], showPods, showUtil)
	}

	w.Flush()
}

func printHeaders(w *tabwriter.Writer, cm *clusterMetric, showPods bool, showUtil bool) {
	if showPods && showUtil {
		fmt.Fprintln(w, "NODE\t NAMESPACE\t POD\t CPU REQUESTS \t CPU LIMITS \t CPU UTIL \t MEMORY REQUESTS \t MEMORY LIMITS \t MEMORY UTIL")

		if len(cm.nodeMetrics) > 1 {
			fmt.Fprintf(w, "* \t *\t *\t %s \t %s \t %s \t %s \t %s \t %s \n",
				cm.cpu.requestString(),
				cm.cpu.limitString(),
				cm.cpu.utilString(),
				cm.memory.requestString(),
				cm.memory.limitString(),
				cm.memory.utilString())

			fmt.Fprintln(w, "\t\t\t\t\t\t\t\t")
		}
	} else if showPods {
		fmt.Fprintln(w, "NODE\t NAMESPACE\t POD\t CPU REQUESTS \t CPU LIMITS \t MEMORY REQUESTS \t MEMORY LIMITS")

		fmt.Fprintf(w, "* \t *\t *\t %s \t %s \t %s \t %s \n",
			cm.cpu.requestString(),
			cm.cpu.limitString(),
			cm.memory.requestString(),
			cm.memory.limitString())

		fmt.Fprintln(w, "\t\t\t\t\t\t")

	} else if showUtil {
		fmt.Fprintln(w, "NODE\t CPU REQUESTS \t CPU LIMITS \t CPU UTIL \t MEMORY REQUESTS \t MEMORY LIMITS \t MEMORY UTIL")

		fmt.Fprintf(w, "* \t %s \t %s \t %s \t %s \t %s \t %s \n",
			cm.cpu.requestString(),
			cm.cpu.limitString(),
			cm.cpu.utilString(),
			cm.memory.requestString(),
			cm.memory.limitString(),
			cm.memory.utilString())

	} else {
		fmt.Fprintln(w, "NODE\t CPU REQUESTS \t CPU LIMITS \t MEMORY REQUESTS \t MEMORY LIMITS")

		if len(cm.nodeMetrics) > 1 {
			fmt.Fprintf(w, "* \t %s \t %s \t %s \t %s \n",
				cm.cpu.requestString(), cm.cpu.limitString(),
				cm.memory.requestString(), cm.memory.limitString())
		}
	}
}

func printNode(w *tabwriter.Writer, name string, nm *nodeMetric, showPods bool, showUtil bool) {
	podNames := make([]string, len(nm.podMetrics))

	i := 0
	for name := range nm.podMetrics {
		podNames[i] = name
		i++
	}
	sort.Strings(podNames)

	if showPods && showUtil {
		fmt.Fprintf(w, "%s \t *\t *\t %s \t %s \t %s \t %s \t %s \t %s \n",
			name,
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.cpu.utilString(),
			nm.memory.requestString(),
			nm.memory.limitString(),
			nm.memory.utilString())

		for _, podName := range podNames {
			pm := nm.podMetrics[podName]
			fmt.Fprintf(w, "%s \t %s \t %s \t %s \t %s \t %s \t %s \t %s \t %s \n",
				name,
				pm.namespace,
				pm.name,
				pm.cpu.requestString(),
				pm.cpu.limitString(),
				pm.cpu.utilString(),
				pm.memory.requestString(),
				pm.memory.limitString(),
				pm.memory.utilString())
		}

		fmt.Fprintln(w, "\t\t\t\t\t\t\t\t")

	} else if showPods {
		fmt.Fprintf(w, "%s \t *\t *\t %s \t %s \t %s \t %s \n",
			name,
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.memory.requestString(),
			nm.memory.limitString())

		for _, podName := range podNames {
			pm := nm.podMetrics[podName]
			fmt.Fprintf(w, "%s \t %s \t %s \t %s \t %s \t %s \t %s \n",
				name,
				pm.namespace,
				pm.name,
				pm.cpu.requestString(),
				pm.cpu.limitString(),
				pm.memory.requestString(),
				pm.memory.limitString())
		}

		fmt.Fprintln(w, "\t\t\t\t\t\t")

	} else if showUtil {
		fmt.Fprintf(w, "%s \t %s \t %s \t %s \t %s \t %s \t %s \n",
			name,
			nm.cpu.requestString(),
			nm.cpu.limitString(),
			nm.cpu.utilString(),
			nm.memory.requestString(),
			nm.memory.limitString(),
			nm.memory.utilString())

	} else {
		fmt.Fprintf(w, "%s \t %s \t %s \t %s \t %s \n", name,
			nm.cpu.requestString(), nm.cpu.limitString(),
			nm.memory.requestString(), nm.memory.limitString())
	}
}
