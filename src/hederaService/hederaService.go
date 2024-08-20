package hederaService

import (
    "context"
    "log"
    "backend/src/config"
    "github.com/hashgraph/hedera-sdk-go/v2"
)

type HederaService struct {
    client *hedera.Client
}

func NewHederaService(config config.Config) *HederaService {
    client := hedera.ClientForName(config.Network)
    client.SetOperator(config.AccountID, config.PrivateKey)

    return &HederaService{client: client}
}

func (hs *HederaService) CreateTopic() (hedera.TopicID, error) {
    txResponse, err := hedera.NewTopicCreateTransaction().
        SetAdminKey(hs.client.GetOperatorPublicKey()).
        Execute(hs.client)
    if err != nil {
        return hedera.TopicID{}, err
    }

    receipt, err := txResponse.GetReceipt(context.Background())
    if err != nil {
        return hedera.TopicID{}, err
    }

    topicID := *receipt.TopicID
    log.Printf("Topic created with ID: %s\n", topicID.String())

    return topicID, nil
}


