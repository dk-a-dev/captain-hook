from datetime import datetime
import json
from flask import Flask
import os
import uuid
import redis
import random

def getPayment():
    """Simulate a payment processing function."""
    return{
        'url':os.getenv('WEBHOOK_URL', ""),
        'webhookId':uuid.uuid4().hex,
        'data':{
            'id':uuid.uuid4().hex,
            'payment': f"PY-{''.join((random.choice('abcdxyzpqr').capitalize() for i in range(5)))}",
            'event': random.choice(["accepted", "completed", "canceled"]),
            'createdAt': datetime.now().strftime("%d/%m/%Y, %H:%M:%S"),
        }
    }

redisAddress= os.getenv('REDIS_ADDRESS')
host,port = redisAddress.split(':')
redisConnection=redis.StrictRedis(host=host, port=port)

app = Flask(__name__)

@app.route('/payment', methods=['POST'])
def payment():
    webhookPayloadJson=json.dumps(getPayment())
    redisConnection.publish('payments',webhookPayloadJson)
    return webhookPayloadJson

if __name__=='__main__':
    app.run(host='0.0.0.0',port=8000)
