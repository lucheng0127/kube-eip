#!/bin/sh
echo 1

/eip_agent --internal-net "$AGENT_SVC_NET" --internal-net "$AGENT_POD_NET" --gateway-ip "$AGENT_EIP_GW_IP" --gateway-dev "$AGENT_EIP_GW_DEV" --log-level "$AGENT_LOG_LEVEL" --eip-cidr "$AGENT_EIP_NET"