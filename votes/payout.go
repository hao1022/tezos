package main

import (
    "fmt"
    "time"
    "../../../model/tezos"
    "strconv"
    "os"
    "strings"
    "bufio"
    "io/ioutil"
)

type ConfigType struct {
    TezosClientPath   string
    Baker             string
    Account           string
    CycleLength       int
    SnapshotInterval  int
    Delegate          string
    DelegateName      string
    FeePercent        int
    StartingCycle     int
    Endpoint          string
    PayoutRecords     string
    PasswordFile      string
    Password          string
}

var Config ConfigType = ConfigType{
    "/home/ubuntu/tezos/tezos-client", // path to tezos-client
    "tz1awXW7wuXy21c66vBudMXQVAPgRnqqwgTH", // baker account
    "infstones",  // key name
    4096, // cycle length, const
    256,  // snapshot interval, const
    "tz1awXW7wuXy21c66vBudMXQVAPgRnqqwgTH", // delegate account
    "infstones", // delegate name
    15, // fee percent, 15% by default
    81, // starting cycle
    "54.188.118.102", // Tezos node to connect to
    "/home/ubuntu/tezos/.payout_records", // payout record file
    "/home/ubuntu/tezos/.password", // password file
    ""} // password, to be input

type StolenBlockType struct {
    Level    int
    Hash     string
    Priority int
    Reward   int
    Fees     int
}

type CycleRewardType struct {
    Realized            int
    Paid                int
    RealizedDifference  int
    EstimatedDifference int
}

type BakerRewardType struct {
    SelfReward int
    FeeReward int
    TotalReward int
}

type DelegatorPayoutType struct {
    Balance             int
    EstimatedRewards    int
    FinalRewards        int
    PayoutOperationHash string
}

func GetEstimatesForCycle(config ConfigType, cycle int) RewardType {
    cycle_length := config.CycleLength
    snapshot_interval := config.SnapshotInterval
    baker := config.Baker
    fee_percent := config.FeePercent

    estimated_rewards := EstimatedRewards(cycle_length, snapshot_interval, baker)
    rewards := CalculateRewardsFor(cycle_length, snapshot_interval, cycle, baker,
                                   estimated_rewards, fee_percent)
    return rewards
}

func GetEstimates(config ConfigType) []RewardType {
    var rewards []RewardType
    current_level := tezos.CurrentLevel()
    current_cycle := current_level.Cycle
    known_cycle := current_cycle + 5
    for cycle := current_cycle; cycle <= known_cycle; cycle++ {
        rewards = append(rewards, GetEstimatesForCycle(config, cycle))
    }
    return rewards
}

func GetActualsForCycle(config ConfigType, cycle int) RewardType {
    cycle_length := config.CycleLength
    snapshot_interval := config.SnapshotInterval
    baker := config.Baker
    fee_percent := config.FeePercent

    var reward RewardType

    hash := HashToQuery(cycle + 1, cycle_length)
    frozen_balance_by_cycle := tezos.FrozenBalanceByCycle(hash, baker)
    for _, balance := range frozen_balance_by_cycle {
        if balance.Cycle == cycle {
	    fee_rewards, _ := strconv.Atoi(balance.Fees)
	    balance_rewards, _ := strconv.Atoi(balance.Rewards)
	    realized_rewards := fee_rewards + balance_rewards
            reward = CalculateRewardsFor(cycle_length, snapshot_interval, cycle,
	                                 baker, realized_rewards, fee_percent)
	}
    }
    return reward
}

func GetPaidCycle() int {
    paid_cycle := Config.StartingCycle - 1
    file, err := os.Open(Config.PayoutRecords)
    if err != nil {
        fmt.Println("An error occured: ", err)
	return 0
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        text := scanner.Text()
	texts := strings.Split(text, ":")

	field_str := texts[0]
	if field_str == "Cycle" {
	    cycle_str := strings.Trim(texts[1], "\n")
            cycle, _ := strconv.Atoi(cycle_str)
	    if cycle > paid_cycle {
                paid_cycle = cycle
            }
	}
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("An error occured: ", err)
	return 0
    }
    return paid_cycle
}

func GetActuals(config ConfigType) []RewardType {
    var rewards []RewardType
    current_level := tezos.CurrentLevel()
    current_cycle := current_level.Cycle
    paid_cycle := GetPaidCycle()
    if paid_cycle == 0 {
	 fmt.Println("Can not get paid cycle")
	 return nil
    }
    delivered_cycle := current_cycle - 6
    for cycle := paid_cycle + 1; cycle <= delivered_cycle; cycle++ {
        rewards = append(rewards, GetActualsForCycle(config, cycle))
    }
    return rewards
}

func PrintReward(reward RewardType) {
    fmt.Printf("Paid Cycle: %d\n", reward.Cycle)
    fmt.Printf("Baker Self Reward: %d\n", reward.BakerRewards.SelfReward)
    fmt.Printf("Baker Fee Reward: %d\n", reward.BakerRewards.FeeReward)
    fmt.Printf("Baker Total Reward: %d\n", reward.BakerRewards.TotalReward)
    for i, _ := range reward.Delegators {
            fmt.Println(reward.Delegators[i])
            fmt.Printf("  Balance: %d\n", reward.DelegatorBalances[i])
            fmt.Printf("  RawReward: %d\n", reward.DelegatorRawRewards[i])
            fmt.Printf("  Reward: %d\n", reward.DelegatorRewards[i])
            fmt.Printf("  Share: %.2f%%\n", reward.DelegatorShares[i])
    }
    fmt.Printf("Staking Balance: %d\n", reward.StakingBalance)
    fmt.Printf("Total Reward: %d\n", reward.TotalReward)
}

func WriteOutPayout(reward RewardType) {
    file, _ := os.OpenFile(Config.PayoutRecords, os.O_APPEND|os.O_WRONLY, 0644)
    defer file.Close()
    file.WriteString(fmt.Sprintf("Cycle:%d\n", reward.Cycle))
    file.WriteString(fmt.Sprintf("Baker Self Reward:%d\n", reward.BakerRewards.SelfReward))
    file.WriteString(fmt.Sprintf("Baker Fee Reward:%d\n", reward.BakerRewards.FeeReward))
    file.WriteString(fmt.Sprintf("Baker Total Reward:%d\n", reward.BakerRewards.TotalReward))
    for i, _ := range reward.Delegators {
        file.WriteString(reward.Delegators[i])
        file.WriteString(fmt.Sprintf("  Balance:%d\n", reward.DelegatorBalances[i]))
        file.WriteString(fmt.Sprintf("  RawReward:%d\n", reward.DelegatorRawRewards[i]))
        file.WriteString(fmt.Sprintf("  Reward:%d\n", reward.DelegatorRewards[i]))
        file.WriteString(fmt.Sprintf("  Share:%.2f%%\n", reward.DelegatorShares[i]))
    }
    file.WriteString(fmt.Sprintf("Staking Balance:%d\n", reward.StakingBalance))
    file.WriteString(fmt.Sprintf("Total Reward:%d\n", reward.TotalReward))
}

func Payout(rewards []RewardType) {
    for _, reward := range rewards {
	PrintReward(reward)
        for i, _ := range reward.Delegators {
	    // format command
	    if reward.DelegatorRewards[i] < 1000 {
	        continue
	    }

	    amount := reward.DelegatorRewards[i]
	    amount_str := strconv.Itoa(amount)

	    err := Transfer(Config, /*account=*/Config.Account,
	                    /*amount=*/amount_str, /*from=*/Config.Baker,
			    /*to=*/reward.Delegators[i])

	    if err == nil {
	        // wait for the prev transaction to get posted
                time.Sleep(70 * time.Second)
	    } else {
		transfer_cmd := fmt.Sprintf("%s -A %s transfer %.3f from %s to %s",
                           /*client=*/Config.TezosClientPath,
			   /*endpoint=*/Config.Endpoint,
			   /*amount=*/float64(amount) / 1000000.0,
			   /*account=*/Config.Account,
			   /*dest=*/reward.Delegators[i])
		fmt.Println("Transfer failed, try manual transfer with ",
		            transfer_cmd)
	    }
        }
        WriteOutPayout(reward)
    }
}

func ReadPassword() {
    dat, err := ioutil.ReadFile(Config.PasswordFile)
    if err != nil {
        fmt.Println("An error occured: ", err)
	return
    }
    Config.Password = strings.TrimSpace(string(dat))
}

func main() {
    tezos.Initialize()

    // read in password
    ReadPassword()
    fmt.Println("Password read")

    for true {
        rewards := GetActuals(Config)
        Payout(rewards)
        time.Sleep(10 * time.Second)
    }
}
