package dbus

func (di *DbusInstance) CreatePubsubChannel(Owner ResourceId, OwnerType ResourceType,
	ResourceIdentifier string) ChannelId {

	n := PubsubChannel{}
	n.ChannelId = GetChannelId()
	n.Owner = Owner
	n.OwnerType = OwnerType
	n.ResourceIdentifier = ResourceIdentifier
	n.Subscribers = make([]PubsubSubscriber, 0)

	di.PubsubChannels[n.ChannelId] = &n
	di.ResourceRegister(n.Owner, n.OwnerType)

	return n.ChannelId
}

func (di *DbusInstance) AddPubsubChannelSubscriber(chanId ChannelId,
	resourceId ResourceId, resourceType ResourceType, channelIn chan []byte) {

	psc := di.PubsubChannels[chanId] //pubsub channel
	ns := PubsubSubscriber{}         //new subscriber
	ns.SubscriberId = resourceId
	ns.SubscriberType = resourceType
	ns.Channel = channelIn

	psc.Subscribers = append(psc.Subscribers, ns)
}

func (di *DbusInstance) PublishTo(chanId uint32, msg []byte) {
	//println("<dbus/pubsub>.PublishTo()", id)
	id := ChannelId(chanId)

	channel := di.PubsubChannels[id]
	di.prefixMessageWithChanId(id, &msg)

	//FIXME? fix non-determinism?
	for _, sub := range channel.Subscribers {
		sub.Channel <- msg
	}
}

func (di *DbusInstance) RemoveChannel(chanId ChannelId) {
	//println("len of di.PubsubChannels:", len(di.PubsubChannels))
	delete(di.PubsubChannels, chanId)
	//println("len of di.PubsubChannels:", len(di.PubsubChannels))
}

//
//
//private
//
//

func (di *DbusInstance) prefixMessageWithChanId(id ChannelId, msg *[]byte) {
	prefix := make([]byte, 4)
	prefix[0] = (uint8)((id & 0x000000ff) >> 0)
	prefix[1] = (uint8)((id & 0x0000ff00) >> 8)
	prefix[2] = (uint8)((id & 0x00ff0000) >> 16)
	prefix[3] = (uint8)((id & 0xff000000) >> 24)
	(*msg) = append(prefix, *msg...)
}
