/**
 * Author: BaldKael
 * Email:  baldkael@sina.com
 * Date:   2020/8/15 16:31
 */

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/BaldKael/MemeWithText/meme"
)

func main() {
	rand.Seed(time.Now().Unix())
	meme := meme.New("source/pics/2.jpg", "target/target.png", "WenQuanYiZenHeiMono", "哦？", -1, 20, 10, []int{0, 0, 139})
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	meme.Resize().AddText().Save().Clean()
	log.Println("done")
}
