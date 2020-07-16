package main

import (
	"fmt"
	"time"
)

func main() {
	obj := Constructor()

	userId := 11556
	tweetId := 1155601
	obj.PostTweet(userId, tweetId)

	feed := obj.GetNewsFeed(userId)
	fmt.Println(feed)

	userId2 := 11557
	obj.Follow(userId2, userId)
	feed2 := obj.GetNewsFeed(userId2)
	fmt.Println(feed2)

	obj.Unfollow(userId2, userId)
	feed3 := obj.GetNewsFeed(userId2)
	fmt.Println(feed3)
}

type Twitter struct {
	users map[int]*User
	// tweets用来存放用户发表的推文；
	tweets map[int]*Tweet
}

type Tweet struct {
	ID   int
	UID  int
	time int64
}

type User struct {
	ID int

	// fans用来存放用户的粉丝（关注者）列表
	followers map[int]*User
	followees map[int]*User // 关注的人,推模式使用，这里不使用

	// feeds用来存放每个用户可以看到的动态
	feed   []*Tweet
	tweets []*Tweet
}

func NewUser(userId int) *User {
	return &User{
		ID:        userId,
		followers: make(map[int]*User),
		followees: make(map[int]*User),
	}
}

func Constructor() Twitter {
	return Twitter{
		users:  make(map[int]*User),
		tweets: make(map[int]*Tweet),
	}
}

// 创建一条新推文
func (this *Twitter) PostTweet(userId int, tweetId int) {
	// 用户是否存在，不存在则创建
	u, ok := this.users[userId]
	if !ok {
		u = NewUser(userId)
		this.users[userId] = u
	}

	// 创建推文
	t := &Tweet{
		ID:   tweetId,
		UID:  userId,
		time: time.Now().UnixNano(),
	}
	this.tweets[tweetId] = t

	u.feed = append(u.feed, t)
	u.tweets = append(u.tweets, t)

	// 给关注了该用户的所有用户的feed中添加这条新创建的 tweet
	for _, f := range u.followers {
		if f != nil {
			f.feed = append(f.feed, t)
		}
	}
}

// 检索最近的十条推文，每个推文都必须由当前用户关注的人或者用户自己发出的。
// 推文需要按时间顺序由最近的开始排序。
func (this *Twitter) GetNewsFeed(userId int) []int {
	var r []int

	u, ok := this.users[userId]
	if !ok {
		return r
	}

	for i := 0; i < len(u.feed); i++ {
		if i == 10 {
			break
		}
		// 倒序取出tweetId
		r = append(r, u.feed[len(u.feed)-1-i].ID)
	}
	return r
}

func (this *Twitter) Follow(ferId int, feeId int) {
	// 不允许自己关注自己
	if ferId == feeId {
		return
	}

	ferUser, ok := this.users[ferId]
	if !ok {
		ferUser = NewUser(ferId)
		this.users[ferId] = ferUser
	}

	feeUser, ok := this.users[feeId]
	if !ok {
		feeUser = NewUser(feeId)
		this.users[feeId] = feeUser
	}

	if _, ok := feeUser.followers[ferId]; ok {
		return
	}

	feeUser.followers[ferId] = ferUser
	ferUser.feed = this.Merge(ferUser.feed, feeUser.feed)
	return
}

func (this *Twitter) Merge(feed1, feed2 []*Tweet) []*Tweet {
	if len(feed1) == 0 {
		return feed2
	}

	if len(feed2) == 0 {
		return feed1
	}

	var f []*Tweet
	i := 0
	j := 0

	// 归并排序
	for i < len(feed1) && j < len(feed2) {
		if feed1[i].time < feed2[j].time {
			f = append(f, feed1[i])
			i++
		} else {
			f = append(f, feed2[j])
			j++
		}
	}

	// 上下的直接加入尾部即可，先前的操作已经将其中一个倒空
	for ; i < len(feed1); i++ {
		f = append(f, feed1[i])
	}

	for ; j < len(feed2); j++ {
		f = append(f, feed2[j])
	}

	return f
}

func (this *Twitter) Unfollow(ferId, feeId int) {
	if ferId == feeId {
		return
	}

	ferUser, ok := this.users[ferId]
	if !ok {
		return
	}

	feeUser, ok := this.users[feeId]
	if !ok {
		return
	}

	if _, ok := feeUser.followers[ferId]; !ok {
		return
	}

	delete(feeUser.followers, ferId)
	var feedNew []*Tweet
	// 将关于feeId相关用的推特 从ferId用户的推荐列表里清除
	for i := range ferUser.feed {
		if ferUser.feed[i].UID != feeUser.ID {
			feedNew = append(feedNew, ferUser.feed[i])
		}
	}
	ferUser.feed = feedNew
}
