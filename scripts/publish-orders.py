import requests
import uuid
import random

item_ids = []
for i in range(1000):
    item_ids.append(uuid.uuid4().__str__())

customers = []
for i in range(100):
    customers.append(uuid.uuid4().__str__())

for i in range(120):
    customer = random.choice(customers)

    num_line_items = random.randint(1, 10)

    line_items = []
    for j in range(num_line_items):
        item_id = random.choice(item_ids)
        line_items.append(
            {
                "item_id": item_id,
                "quantity": random.randint(1, 10),
                "price": random.randint(1, 10000),
            }
        )

    order = {
        "customer_id": customer,
        "line_items": line_items,
    }

    # If your port is different, change this URL
    r = requests.post("http://localhost:3000/orders", json=order)
    r.status_code
    print("posted order", i + 1)
