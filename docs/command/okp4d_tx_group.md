## okp4d tx group

Group transaction subcommands

```
okp4d tx group [flags]
```

### Options

```
  -h, --help   help for group
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID (default "okp4d")
```

### SEE ALSO

* [okp4d tx](okp4d_tx.md)	 - Transactions subcommands
* [okp4d tx group create-group](okp4d_tx_group_create-group.md)	 - Create a group which is an aggregation of member accounts with associated weights and an administrator account.
* [okp4d tx group create-group-policy](okp4d_tx_group_create-group-policy.md)	 - Create a group policy which is an account associated with a group and a decision policy. Note, the '--from' flag is ignored as it is implied from [admin].
* [okp4d tx group create-group-with-policy](okp4d_tx_group_create-group-with-policy.md)	 - Create a group with policy which is an aggregation of member accounts with associated weights, an administrator account and decision policy.
* [okp4d tx group draft-proposal](okp4d_tx_group_draft-proposal.md)	 - Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).
* [okp4d tx group exec](okp4d_tx_group_exec.md)	 - Execute a proposal
* [okp4d tx group leave-group](okp4d_tx_group_leave-group.md)	 - Remove member from the group
* [okp4d tx group submit-proposal](okp4d_tx_group_submit-proposal.md)	 - Submit a new proposal
* [okp4d tx group update-group-admin](okp4d_tx_group_update-group-admin.md)	 - Update a group's admin
* [okp4d tx group update-group-members](okp4d_tx_group_update-group-members.md)	 - Update a group's members. Set a member's weight to "0" to delete it.
* [okp4d tx group update-group-metadata](okp4d_tx_group_update-group-metadata.md)	 - Update a group's metadata
* [okp4d tx group update-group-policy-admin](okp4d_tx_group_update-group-policy-admin.md)	 - Update a group policy admin
* [okp4d tx group update-group-policy-decision-policy](okp4d_tx_group_update-group-policy-decision-policy.md)	 - Update a group policy's decision policy
* [okp4d tx group update-group-policy-metadata](okp4d_tx_group_update-group-policy-metadata.md)	 - Update a group policy metadata
* [okp4d tx group vote](okp4d_tx_group_vote.md)	 - Vote on a proposal
* [okp4d tx group withdraw-proposal](okp4d_tx_group_withdraw-proposal.md)	 - Withdraw a submitted proposal

