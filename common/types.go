package common

//BlockHeaderType type
type BlockHeaderType struct{
	Protocol		 string   `json:"protocol"`
	ChainID			 string   `json:"chain_id"`
	Hash			 string   `json:"hash"`
	Level			 int	  `json:"level"`
	Proto			 int	  `json:"proto"`
	Predecessor		 string   `json:"predecessor"`
	Timestamp		 string   `json:"timestamp"`
	ValidationPass	 int	  `json:"validation_pass"`
	OperationsHash	 string   `json:"operations_hash"`
	Fitness			 []string `json:"fitness"`
	Context			 string   `json:"context"`
	Priority		 int	  `json:"priority"`
	Nonce			 string   `json:"proof_of_work_nonce"`
	Signature		 string   `json:"signature"`
}

//BalanceUpdateType type
type BalanceUpdateType struct{
	Kind	 string   `json:"kind"`
	Category string   `json:"category"`
	Contract string   `json:"contract"`
	Delegate string   `json:"delegate"`
	Level	 int	  `json:"level"`
	Change	 string   `json:"change"`
}

//BlockMetadataType type
type BlockMetadataType struct{
	Protocol	   string				  `json:"protocol"`
	Baker		   string				  `json:"baker"`
	BalanceUpdates []BalanceUpdateType	  `json:"balance_updates"`
}

//EndorsingRightType type
type EndorsingRightType struct{
	Level		  int	   `json:"level"`
	Delegate	  string   `json:"delegate"`
	Slots		  []int    `json:"slots"`
	EstimatedTime int	   `json:"estimated_time"`
}

//BakingRightType type
type BakingRightType struct{
	Level		  int	   `json:"level"`
	Delegate	  string   `json:"delegate"`
	Priority	  int	   `json:"priority"`
	EstimatedTime int	   `json:"estimated_time"`
}

//CurrentLevelType type
type CurrentLevelType struct{
	Level				 int   `json:"level"`
	LevelPosition		 int   `json:"level_position"`
	Cycle				 int   `json:"cycle"`
	CyclePosition		 int   `json:"cycle_position"`
	VotingPeriod		 int   `json:"voting_period"`
	VotingPeriodPosition int   `json:"voting_period_position"`
	ExpectedCommitment	 bool  `json:"expected_commitment"`
}

//OperationResultType type
type OperationResultType struct{
	Status			string	`json:"status"`
	BalanceUpdates	[]BalanceUpdateType    `json:"balance_updates"`
}

//OperationMetadataType type
type OperationMetadataType struct{
	BalanceUpdates []BalanceUpdateType	`json:"balance_update"`
	Result		   OperationResultType	`json:"operation_result"`
}

//OperationContentType type
type OperationContentType struct{
	Kind		string	 `json:"kind"`
	Source		string	 `json:"source"`
	Fee			string	 `json:"fee"`	  // mutez
	Counter		string	 `json:"counter"`
	GasLimit	string	 `json:"gas_limit"`
	StorageLimit string  `json:"storage_limit"`
	Destination string	 `json:"destination"`
	Amount		string	 `json:"amount"`  // mutez
}

//OperationContentAndResultType type
type OperationContentAndResultType struct{
	OperationContentType
	Metadata	OperationMetadataType `json:"metadata"`
}

//OperationContentsAndResultsType type
type OperationContentsAndResultsType struct{
	Contents	[]OperationContentAndResultType `json:"contents"`
}

//OperationType type
type OperationType struct{
	Protocol string   `json:"protocol"`
	Hash	 string   `json:"hash"`
	Branch	 string   `json:"branch"`
	Contents []OperationContentType  `json:"contents"`
}

//PreapplyResultType type
type PreapplyResultType struct{
	Kind	 string   `json:"kind"`
	ID		 string   `json:"id"`
	Implicit string   `json:"implicit"`
}

//type OperationResultType struct{
//	  Contents	  []OperationContentType  `json:"contents"`
//	  Signature   string				  `json:"signature"`
//}

//FrozenBalanceByCycleType type
type FrozenBalanceByCycleType struct{
	Cycle	int		 `json:"cycle"`
	Deposit string	 `json:"deposit"`
	Fees	string	 `json:"fees"`
	Rewards string	 `json:"rewards"`
}

//CycleInfoType type
type CycleInfoType struct{
	LastRoll   []string `json:"last_roll"`
	Nonces	   []string `json:"nonces"`
	RandomSeed string	`json:"random_seed"`
	Snapshot   int		`json:"roll_snapshot"`
}
