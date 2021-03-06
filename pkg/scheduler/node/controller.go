//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package node

import (
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/scheduler/context"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type NodeController struct {
	context *context.Context
	nodes   chan *types.Node
	active  bool
}

func (nc *NodeController) Watch(node chan *types.Node) {

	var (
		stg = nc.context.GetStorage()
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case n := <-nc.nodes:
				{
					if !nc.active {
						log.Debug("NodeController: skip management cause it is in slave mode")
						continue
					}

					log.Debugf("Node check state: %s", n.Meta.Name)
					if n.Alive {
						log.Debugf("Node set alive, try to provision on it pods: %s", n.Meta.Name)
						node <- n
						continue
					}

					log.Debugf("Node set offline, try to move all pods to another")

				}
			}
		}
	}()

	stg.Node().Watch(nc.context.Background(), nc.nodes)
}

func (nc *NodeController) Pause() {
	nc.active = false
}

func (nc *NodeController) Resume() {

	nc.active = true
	log.Debug("NodeController: start check pods state")
}

func NewNodeController(ctx *context.Context) *NodeController {
	sc := new(NodeController)
	sc.context = ctx
	sc.active = false
	sc.nodes = make(chan *types.Node)
	return sc
}
