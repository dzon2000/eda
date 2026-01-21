# Putting everything together

Now as we know the basics, let's build the real system. It will utilize all the knowledge and implement
simple ordering service:

- Creating the order
- Consuming order event and proceeding with the payment
- Consuming payment event and fulfilling the order

## Order service

Order Service:
- validates commands
- persists orders
- emits OrderCreated (outbox publisher)
- does **NOT** handle payments or fulfillment

