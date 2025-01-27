package config

import (
	"log"

	conversationsModels "Orbyters/models/conversations"
	userModels "Orbyters/models/users"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"

	"gorm.io/gorm"
)

func ApplySeeds(db *gorm.DB) {
	seedRoles(db)
	seedMessageTypes(db)
	seedSubscriptions(db)
}

func seedRoles(db *gorm.DB) {
	roles := []userModels.Role{
		{Name: "Admin"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		var existingRole userModels.Role
		err := db.Where("name = ?", role.Name).First(&existingRole).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error creating role: %v", err)
		}

		if existingRole.Id == 0 {
			existingRole.Name = role.Name
			if err := existingRole.CreateRole(db); err != nil {
				log.Fatalf("Error saving role: %v", err)
			}
			log.Printf("Role '%s' created.", role.Name)
		} else {
			log.Printf("Role '%s' already existing.", role.Name)
		}
	}
}

func seedMessageTypes(db *gorm.DB) {
	messageTypes := []conversationsModels.MessageType{
		{Type: "user"},
		{Type: "assistant"},
	}

	for _, messageType := range messageTypes {
		var existingMessageType conversationsModels.MessageType
		err := db.Where("Type = ?", messageType.Type).First(&existingMessageType).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error creating message type: %v", err)
		}

		if existingMessageType.Id == 0 {
			existingMessageType.Type = messageType.Type
			if err := existingMessageType.CreateMessageType(db); err != nil {
				log.Fatalf("Error saving message type: %v", err)
			}
			log.Printf("Message type '%s' created.", messageType.Type)
		} else {
			log.Printf("Message type '%s' already existing.", messageType.Type)
		}
	}
}

func seedSubscriptions(db *gorm.DB) {
	subscriptions := []userModels.Subscription{
		{Price: 120.00, Title: "Moon", Description: "moonDescription", StripeProductId: ""},
	}

	for _, subscription := range subscriptions {
		var existingSubscription userModels.Subscription
		err := db.Where("Title = ?", subscription.Title).First(&existingSubscription).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Fatalf("Error creating subscription: %v", err)
		}

		if existingSubscription.Id == 0 {
			existingSubscription.Title = subscription.Title
			existingSubscription.Price = subscription.Price
			existingSubscription.Description = subscription.Description
			existingSubscription.StripeProductId = subscription.StripeProductId

			productId, err := refreshProduct(subscription.Title, subscription.Price)
			if err != nil {
				log.Fatalf("Error creating product: %v", err)
			}
			existingSubscription.StripeProductId = productId

			if err := existingSubscription.CreateSubscription(db); err != nil {
				log.Fatalf("Error saving subscription: %v", err)
			}
			log.Printf("Subscription '%s' created.", subscription.Title)
		} else {
			log.Printf("Subscription '%s' already existing.", subscription.Title)
		}
	}
}

func refreshProduct(productName string, productPrice float64) (string, error) {
	stripe.Key = StripeKey

	productParams := &stripe.ProductParams{
		Name: stripe.String(productName),
	}
	createdProduct, err := product.New(productParams)
	if err != nil {
		return "", err
	}

	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(string(stripe.CurrencyEUR)),
		UnitAmount: stripe.Int64(int64(productPrice * 100)),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(string(stripe.PriceRecurringIntervalMonth)),
		},
		Product: stripe.String(createdProduct.ID),
	}
	_, err = price.New(priceParams)
	if err != nil {
		return "", err
	}

	return createdProduct.ID, nil
}
