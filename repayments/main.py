from attr import define


PAYMENTS = {
    "andrew": -166.98,
    "austin": -57.15,
    "jix":     52.90,
    "karen": -132.69,
    "kjiao": 474.13,
    "nikevi": -131.00,
    "tina": -88.13,
    "stine": 147.67,
    "megan": 147.67,
    "neil": 147.67,
    "nurr": -184.14,
    "tim":     -1.57,
    "yiwen": -208.36,
}


VENMO = {
    "andrew": "@andrew-fang-6",
    "yiwen": "@yiwen-song",
    "tim": "@timothy-hsu-4",
    "karen": "@karending",
}


@define
class Payment:
    from_: str
    to_: str
    amount: float

    def __str__(self):
        return f'{self.to_} requests {self.from_} ({VENMO.get(self.from_)}) for ${self.amount}'


def main() -> None:
    """Runs repayments"""

    print(sum(PAYMENTS.values()))
    assert(sum(PAYMENTS.values()) <= 1)

    negatives = {name: amt for name, amt in PAYMENTS.items() if amt < 0}
    positives = {name: amt for name, amt in PAYMENTS.items() if amt > 0}

    payments = []

    while negatives:
        payer = next(iter(negatives.keys()))
        payer_amount_left = negatives[payer]

        recipient = next(iter(positives.keys()))
        recipient_amount_left = positives[recipient]

        payment_amount = round(min(-payer_amount_left, recipient_amount_left), 2)
        payments.append(Payment(
            from_=payer,
            to_=recipient,
            amount=payment_amount,
        ))

        new_payer_balance = payer_amount_left + payment_amount
        if new_payer_balance >= -.01:
            negatives.pop(payer)
        else:
            negatives[payer] = new_payer_balance

        new_recipient_balance = recipient_amount_left - payment_amount
        if new_recipient_balance <= .01:
            positives.pop(recipient)
        else:
            positives[recipient] = new_recipient_balance

        print(positives)
        print(negatives)

    for payment in payments:
        print(payment)


if __name__ == "__main__":
    main()
