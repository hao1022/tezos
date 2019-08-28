package main

import (
    "fmt"
    "strconv"
    "errors"
    "encoding/json"
    "../../../model/tezos"

    "strings"
    "bufio"
    "time"
    "os/exec"
    "github.com/kr/pty"
    "github.com/btcsuite/btcutil/base58"
)

func Sign(config ConfigType, operation string, account string) (string, error) {
    signature := ""
    op_bytes := fmt.Sprintf("0x03%s", operation)
    process := exec.Command(
	    config.TezosClientPath, "-A",
	    config.Endpoint, "sign", "bytes", op_bytes,
	    "for", account)

    cmd_str := fmt.Sprintf("%s -A %s sign bytes %s for %s",
                           config.TezosClientPath, config.Endpoint,
                           op_bytes, account)

    tty, _ := pty.Start(process)
    defer tty.Close()

    time.Sleep(1000 * time.Millisecond)

    // redirect tty stdin
    go func() {
        tty.Write([]byte(config.Password + "\n"))
    }()
    time.Sleep(100 * time.Millisecond)
    scanner := bufio.NewScanner(tty)
    for scanner.Scan() {
	line := scanner.Text()
	if strings.HasPrefix(line, "Signature: ") {
	    signature = line[11:]
	    return signature, nil
	}
    }
    return signature, errors.New(cmd_str + " failed")
}

func Transfer(config ConfigType, account string, amount string, from string, to string) error {
    header := tezos.Header()

    counter := tezos.Counter(from)
    counter = counter + 1
    count := strconv.Itoa(counter)

    txn := tezos.OperationContentType{
        Kind: "transaction",
	Amount: amount,
	Source: from,
        Fee: "3000",
	Counter: count,
	GasLimit: "11000",
	StorageLimit: "0",
	Destination: to}

    txn_str, _ := json.Marshal(txn)
    signature := "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q"

    run_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}",
                           header.Hash, txn_str, signature)

    {
        run_json_result, err := tezos.RunOperation(run_json)
        if err != nil {
	    fmt.Println("RunOperation failed: ", err)
	    fmt.Println("Transaction Json: ", run_json)
	    return err
	}
	if run_json_result.Contents[0].Metadata.Result.Status != "applied" {
	    fmt.Println("RunOperation failed: ", run_json)
            return errors.New("RunOperation json not applied")
        }
    }

    sign_json := fmt.Sprintf("{\"branch\": \"%s\", \"contents\": [%s]}", header.Hash, txn_str)
    sign_json_result, forge_err := tezos.ForgeOperations(sign_json)
    if forge_err != nil {
        fmt.Println("ForgeOperations failed: ", forge_err)
	return forge_err
    }

    sig, sign_err := Sign(config, sign_json_result, account)
    if sign_err != nil {
       fmt.Println("Sign failed: ", sign_err)
       return sign_err
    }

    base16sig := base58.Decode(sig)
    truesig := base16sig[5:len(base16sig)-4]

    signed_op := fmt.Sprintf("\"%s%x\"", sign_json_result, truesig)

    {
        preapply_json := fmt.Sprintf(
	        "[{\"protocol\": \"%s\", \"branch\": \"%s\", \"contents\": [%s], \"signature\": \"%s\"}]",
	        header.Protocol, header.Hash, txn_str, sig)
        _, err := tezos.PreapplyOperations(preapply_json)
        if err != nil {
	    fmt.Println("Preapply failed: ", err)
	    return err
        }
    }

    {
        txn_hash, err := tezos.Injection(signed_op)
	if err != nil {
	    fmt.Println("Injection failed: ", err)
	    return err
	}
	fmt.Printf("Transfer Txn Hash: %s\n", txn_hash)
    }

    return nil
}
