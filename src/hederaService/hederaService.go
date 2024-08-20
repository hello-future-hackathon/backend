package hederaService

import (
    "log"
    "backend/src/config"
    "github.com/hashgraph/hedera-sdk-go/v2"
)

type HederaService struct {
    client *hedera.Client
}

func NewHederaService(cfg config.Config) (*HederaService, error) {
    client, err := hedera.ClientForName(cfg.Network)
    if err != nil {
        return nil, err
    }
    
    client.SetOperator(cfg.AccountID, cfg.PrivateKey)
    return &HederaService{client: client}, nil
}

func (hs *HederaService) CreateTopic() (hedera.TopicID, error) {
    txResponse, err := hedera.NewTopicCreateTransaction().
        SetAdminKey(hs.client.GetOperatorPublicKey()).
        Execute(hs.client)
    if err != nil {
        return hedera.TopicID{}, err
    }

    receipt, err := txResponse.GetReceipt(hs.client)
    if err != nil {
        return hedera.TopicID{}, err
    }

    topicID := *receipt.TopicID
    log.Printf("Topic created with ID: %s\n", topicID.String())

    return topicID, nil
}

func (hs *HederaService) SendMessage(topicID hedera.TopicID, message string) (hedera.TransactionResponse, error) {
    txResponse, err := hedera.NewTopicMessageSubmitTransaction().
        SetTopicID(topicID).
        SetMessage([]byte(message)).
        Execute(hs.client)
    if err != nil {
        return hedera.TransactionResponse{}, err
    }

    log.Printf("Message sent to topic %s: %s\n", topicID.String(), message)
    return txResponse, nil
}

func (hs *HederaService) SubscribeToTopic(topicID hedera.TopicID, messageHandler func(hedera.TopicMessage)) error {
    _, err := hedera.NewTopicMessageQuery().
        SetTopicID(topicID).
        Subscribe(hs.client, messageHandler)
    if err != nil {
        return err
    }

    log.Printf("Subscribed to topic: %s\n", topicID.String())
    return nil
}
