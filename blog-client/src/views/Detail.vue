<template>
  <canvas-background></canvas-background>
    <div style="position: fixed; top: 0; height: 50px; width: 100%; z-index: 999; "><TopBar></TopBar></div>

    <div class="tabs">
        <n-card>
            <n-h1>{{articleInfo.title}}</n-h1>
            <div style="height: 75px; background-color: #FCFAF7;">
                <n-avatar @click="toOtherUser" round size="medium" :src=userUrl style="position: relative; left: 20px; top: 20px; cursor: pointer;"/>
                <text style="position: relative; left: 36px; color: #808080;">发布时间：{{articleInfo.createdAt}} </text>
                <div style="position: relative; left: 70px; color: #808080;">
                    文章分类：
                    <n-tag type="warning">{{categoryName}}</n-tag>
                </div>
                <n-button v-if="self" @click="toUpdate" ghost style="bottom: 45px; left: 805px;" color="#7B3DE0">修改</n-button>
                <n-button v-if="self" @click="toDelete" ghost style="bottom: 45px; left: 815px;" color="#7B3DE0">删除</n-button>
            </div>
            <n-divider />
            <n-button v-if=!liked text @click="newLike" style="position: absolute; right: 120px; top: 25px; cursor: pointer;" round ghost color="red">
              <template #icon>
                <n-icon>
                  <thumbs-up-outline />
                </n-icon>
              </template>
              点赞
            </n-button>
            <n-button  v-else text @click="unLike" style="position: absolute; right: 120px; top: 25px; cursor: pointer;" round ghost color="red">
              <template #icon>
                <n-icon>
                  <thumbs-up />
                </n-icon>
              </template>
              取消点赞
            </n-button>
            <n-button v-if=!collected text @click="newCollect" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#FFA876">
                <template #icon>
                    <n-icon>
                        <star-outline />
                    </n-icon>
                </template>
                收藏
            </n-button>
            <n-button v-else text @click="unCollect" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="#FFA876">
                <template #icon>
                    <n-icon>
                        <star />
                    </n-icon>
                </template>
                已收藏
            </n-button>
            <div class="article-content">
                <div v-html="articleInfo.content"></div>
            </div>
        </n-card>
        <n-card style="padding: 15px;">
          <div class="comment-content">
            <n-text style="font-weight: bold;">评论 :</n-text>
            <n-button text @click="addComment" style="position: absolute; right: 40px; top: 25px; cursor: pointer;" round ghost color="green">
              <template #icon>
                <n-icon>
                  <chatbox-ellipses />
                </n-icon>
              </template>
              添加评论
            </n-button>
          </div>
          <div class="comment-list" style="margin-top: 40px;width:100%;">
            <div v-for="(comment, index) in commentList" style="margin-bottom:15px;border-top: 1px solid #ccc;padding-top:10px;">
              <div class="comment-user-avatar">
                <n-avatar :src="comment.user_avatar" style="width:40px;height:40px;float:left;" />
                <span class="comment-user-text">{{comment.user_name}}</span>
                <span class="comment-user-text" style="color:#5a5e66;float: right;"> 发表于：{{ comment.create_time }}</span>
              </div>
              <div class="comment-part" style="padding-left: 60px;">
                <div style="width:94%;display: inline-block;">{{comment.content}}</div>
                <div style="width:4%;display: inline-block;">
                  <n-button text @click="deleteComment(comment.id)" style="cursor: pointer;color: darkred;display:inline-block;float: right;" round ghost >删除</n-button>
                </div>
              </div>
            </div>
            <n-pagination v-if="commentList.length > 0" @update:page="loadComments" v-model:page="pageInfo.pageNum" :page-count="pageInfo.pageCount" />
          </div>
        </n-card>
    </div>
  <n-modal v-model:show="showModal">
    <div style="width: 600px; height: 250px; background: white;">
      <n-card title="添加评论" :bordered="false">
        <div style="width:100%; margin: 0 auto;">
          <n-input v-model:value="commentData.content" type="textarea" placeholder="请输入评论"/>
        </div>
      </n-card>
      <div style="position: absolute; right: 100px; bottom: 30px;">
        <n-button type="default" @click="closeSubmitModal">
          取消
        </n-button>
      </div>
      <div style="position: absolute; right: 30px; bottom: 30px;">
        <n-button type="primary" @click="submitComment">
          提交
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup>
import {ref,reactive,inject, onMounted} from 'vue'
import TopBar from '../components/TopBar.vue'
import { Star, StarOutline,ThumbsUp,ThumbsUpOutline,ChatboxEllipses } from "@vicons/ionicons5"

import {useRouter, useRoute} from 'vue-router'
import CanvasBackground from "../components/CanvasBackground.vue";
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")
const dialog = inject("dialog")

const articleInfo = ref({})
const categoryName = ref("")
const user = ref({})
const userUrl = ref("")
const collected = ref(false)
const index = ref(0)
const liked = ref(false)
const like_id = ref(0)
const self = ref(false)
const showModal = ref(false)

const commentData = reactive({
  article_id: 0,
  content:"",
})

const commentList = ref([])

const pageInfo = reactive({
  pageNum:1,
  pageSize:5,
  pageCount:0,
  count:0,
  articleId:0
})

onMounted(() => {
    loadArticle()
    loadComments()
})

const loadArticle = async() => {
    let res1 = await axios.get("article/" + route.query.id)

    if (res1.data.code == 200) {
        articleInfo.value = res1.data.data.article 
        let res2 = await axios.get("category/" + res1.data.data.article.category_id)
        if (res2.data.code == 200) {
            categoryName.value = res2.data.data.categoryName
        }
        let res3 = await axios.get("user/briefInfo/" + res1.data.data.article.user_id)
        if (res3.data.code == 200) {
            user.value = res3.data.data
            userUrl.value = serverUrl + user.value.avatar
            if (user.value.id == user.value.loginId) {
                self.value = true
            }
        }    
        let res4 = await axios.get("collects/" + route.query.id) 

        if (res4.data.code == 200) {
            collected.value = res4.data.data.collected
            index.value = res4.data.data.index
        }

      let res5 = await axios.get("article/like/" + route.query.id)

      if (res5.data.code == 200) {
        liked.value = res5.data.data.liked
        like_id.value = res5.data.data.like_id
      }

      pageInfo.articleId = route.query.id
      loadComments()
    }
}

const loadComments = async(pageNum = 0) =>{
  if (pageNum != 0){
    pageInfo.pageNum = pageNum;
  }
  let res = await axios.post(`/article/comment_list?pageNum=${pageInfo.pageNum}&pageSize=${pageInfo.pageSize}&articleId=${pageInfo.articleId}`)
  if (res.data.code == 200) {
    let comments = res.data.data.comments
    let commentUser = res.data.data.commentUsers

    for (let index in comments) {
      let comment = comments[index]
      comments[index]["create_time"] = comments[index]["create_time"].replace("T", " ").replace("\+08:00", " ")
      comments[index]["user_name"] = commentUser[comment.reviewer].userName
      comments[index]["user_avatar"] = serverUrl+commentUser[comment.reviewer].avatar
    }

    commentList.value = comments
  }
  pageInfo.count = res.data.data.count;
  pageInfo.pageCount = parseInt(pageInfo.count / pageInfo.pageSize) + (pageInfo.count % pageInfo.pageSize > 0 ? 1 : 0)
}

const newCollect = async() => {
    let res = await axios.put("collects/new/" + route.query.id)

    if (res.data.code == 200) {
        message.warning("已收藏", {showIcon: false})  
        loadArticle()  
    }
}

const unCollect = async() => {
    let res = await axios.delete("collects/" + index.value)

    if (res.data.code == 200) {
        message.warning("取消收藏", {showIcon: false})  
        loadArticle()  
    }
}

const newLike = async() => {
  let res = await axios.put("article/newLike/" + route.query.id)

  if (res.data.code == 200) {
    message.warning("已点赞", {showIcon: false})
    loadArticle()
  }
}

const unLike= async() => {
  let res = await axios.delete("article/unLike/" + like_id.value)

  if (res.data.code == 200) {
    message.warning("取消点赞", {showIcon: false})
    loadArticle()
  }
}

const addComment = () => {
  showModal.value = true
}

const closeSubmitModal = () => {
  showModal.value = false
  commentData.article_id = 0
  commentData.content = ""
}

const submitComment = async() => {
  commentData.article_id = parseInt(route.query.id)
  let res = await axios.post("article/addComment", {
    article_id: commentData.article_id,
    content: commentData.content,
  })

  if (res.data.code == 200) {
    message.warning("提交成功！", {showIcon: false})
    showModal.value = false
    commentData.content = ""
    commentData.article_id = 0
    loadComments()
  }
}

const deleteComment = async(comment_id) => {
  let res = await axios.post("article/delComment/"+comment_id)
  if (res.data.code == 200) {
    message.warning("已删除！", {showIcon: false})
    loadComments()
  }
}

const toOtherUser = () => {
    if (user.value.id == user.value.loginId) {
        router.push({
            path: "/myself",
            query: {
                id: user.value.id
            }
        })   
    } else {
        router.push({
            path: "/others",
            query: {
                id: user.value.id
            }
        })
    }
}

const toUpdate = () => {
    router.push({
        path: "/update",
        query: {
            id: articleInfo.value.id
        }
    })
}

const toDelete = async (blog) => {
    dialog.warning({
      title: '警告',
      content: '是否要删除',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
            let res = await axios.delete("article/" + articleInfo.value.id)
            if(res.data.code == 200){
                message.info(res.data.msg)
                goback()
            }else{
                message.error(res.data.msg)
            }  
        },
        onNegativeClick: () => {}
    })    
}

const goback= () => {
    router.go(-1)    
}

</script>

<style lang="scss" scoped>
.tabs {
    position: absolute;
    top: 75px;
    left: 0;
    right: 0;
    margin: auto;
    width: 1000px;
    height: 88%;
    overflow: hidden;
    overflow-y: auto;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px; 
}
.article-content img{
    max-width: 100% !important;
}
::-webkit-scrollbar {
  width:5px;
}

.comment-user-avatar {
  display: inline-block;
  height: 40px;
  line-height: 40px;
  width: 100%;
}

.comment-user-text {
  display: inline-block;
  height: 40px;
  line-height: 40px;
  margin-left: 20px;
}
</style>