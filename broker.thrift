
service Broker {
	// publish a msg on a topic. If the topic is not subscribed the msg is lost.
	void publish(1: string topic, 2: string message)

	// subscribe on a topic. returns a key representing the subscription, that
	// can be fed into poll() for receiving messages. A subscription expires
	// automatically if not used for 30s. The topic doesn't need to exist yet.
	string subscribe(1: string topic)

	// given a subscription key, poll for messages. yields an error if key is not valid.
	// A subscription expires automatically after 30s if not polled.
	// max_msgs indicates the maximums number of msgs to return for this call at once.
	list<string> poll(1: string key, 2: i32 maxMsgs)

	// Currently active topics and subscriptions.
	list<SubscribedTopic> activeSubscription()
}

struct SubscribedTopic {
	1: required string topic
	2: required i64    subscribers
}
