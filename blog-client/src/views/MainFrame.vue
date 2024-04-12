<template>
  <canvas-background></canvas-background>
    <div>
        <div style="position: fixed; top: 0; height: 60px; width: 100%; z-index: 999; ">
          <div style="position: fixed; top: 0; height: 60px; width: 80%; z-index: 999; ">
            <div class="card">
              <n-popselect @update:value="searchByCategory" v-model:value="selectedCategory" :options="categoryOptions" trigger="click">
                <n-button text style="position:absolute; left: 50px; top: 22px; font-size: 18px;">{{categoryName}}</n-button>
              </n-popselect>

              <n-input v-model:value="pageInfo.keyword" round placeholder="请输入关键字" style="position:absolute; left: 15%; top: 15px; width:70%; background-color: #F3F0F9;" />

              <n-button @click="loadArticles(0)" round color="#7B3DE0" style="position:absolute; right: 20px; top: 15px;">
                <template #icon>
                  <n-icon>
                    <search />
                  </n-icon>
                </template>
                搜索
              </n-button>
            </div>
          </div>
          <TopBar></TopBar>
        </div>

        <div class="tabs">
            <n-card v-if="userStore.token">
                <div v-for="(article,index) in articleList" style="margin-bottom:15px">
                    <n-card v-if="article.head_image" @click="toDetail(article)" style="cursor: pointer;" hoverable >
                        <n-image width="200" :src=serverUrl+article.head_image style="float: left" />
                        <div style="position: absolute; left: 240px; width: 690px;">
                            <text style="font-weight:bold; font-size: 20px;">
                              {{article.title}}
                              <text style="font-weight: normal;"> — {{article.user_name}}</text>
                            </text>
                            <p>{{article.content+"..."}}</p>
                          <div class="article-other-info" style="display: inline-block;width:100%;">
                            <div class="other-info-part">发布时间：{{article.created_at.replace("T", " ").replace("\+08:00", " ")}}</div>
                            <div class="other-info-part" style="color: red;">
                              <n-icon>
                                <thumbs-up />
                              </n-icon>
                              点赞：{{article.like_count}}
                            </div>
                            <div class="other-info-part" style="color: green;">
                              <n-icon>
                                <chatbox-ellipses />
                              </n-icon>
                              评论：{{article.comment_count}}
                            </div>
                          </div>
                        </div>
                    </n-card>
                    <n-card v-else style="cursor: pointer;" hoverable >
                            <div style="height: 140px; ">
                                <text @click="toDetail(article)" style="font-weight:bold; font-size: 20px;">
                                  {{article.title}}
                                  <text style="font-weight: normal;"> — {{article.user_name}}</text>
                                </text>
                                <p @click="toDetail(article)" >{{article.content+"..."}}</p>
                                <div class="article-other-info" style="display: inline-block;width:100%;">
                                  <div class="other-info-part">
                                    发布时间：{{article.created_at.replace("T", " ").replace("\+08:00", " ")}}
                                  </div>
                                  <div class="other-info-part" style="color: red;">
                                    <n-icon>
                                      <thumbs-up />
                                    </n-icon>
                                    点赞：{{article.like_count}}
                                  </div>
                                  <div class="other-info-part" style="color: green;">
                                    <n-icon>
                                      <chatbox-ellipses />
                                    </n-icon>
                                    评论：{{article.comment_count}}
                                  </div>
                                </div>
                            </div>
                        </n-card>
                </div>

                <n-pagination v-if="articleList.length>0" @update:page="loadArticles" v-model:page="pageInfo.pageNum" :page-count="pageInfo.pageCount" />
                <div v-else class="card-text" style="color: red;">未获取到数据！</div>
            </n-card>
            <div v-if="!userStore.token" class="card-text" style="text-decoration: underline;" @click="router.push('/login')">
              暂无数据，请先登录！
            </div>
        </div>
    </div>

</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import {ref,reactive,inject,onMounted,computed} from 'vue'
import {Search} from '@vicons/ionicons5'
import CanvasBackground from "../components/CanvasBackground.vue";
const userStore = UserStore()

import {useRouter, useRoute} from 'vue-router'
import {UserStore} from "../stores/UserStore.js";
import { ThumbsUp,ChatboxEllipses } from "@vicons/ionicons5"
const router = useRouter()
const route = useRoute()

const serverUrl = inject("serverUrl")
const axios = inject("axios")
const message = inject("message")

const selectedCategory = ref(0)
const categoryOptions = ref([])
const articleList = ref([])
const pageInfo = reactive({
  pageNum:1,
  pageSize:5,
  pageCount:0,
  count:0,
  keyword:"",
  categoryId:0
})

onMounted(()=>{
    loadArticles()
    loadCategories()
})

const loadArticles = async(pageNum = 0) =>{
    if (pageNum != 0){
        pageInfo.pageNum = pageNum;
    }
    let res = await axios.post(`/article/list?keyword=${pageInfo.keyword}&pageNum=${pageInfo.pageNum}&pageSize=${pageInfo.pageSize}&categoryId=${pageInfo.categoryId}`)
    // console.log(res)
    if (res.data.code == 200) {
      let articles = res.data.data.article
      let articleCount = res.data.data.articleCount
      let articleUsers = res.data.data.articleUsers

      for (let index in articles) {
        let article = articles[index]
        articles[index]["comment_count"] = articleCount[article.id].comment_count
        articles[index]["like_count"] = articleCount[article.id].like_count
        articles[index]["user_name"] = articleUsers[article.user_id].userName
      }
        articleList.value = articles
    }
    pageInfo.count = res.data.data.count;
    pageInfo.pageCount = parseInt(pageInfo.count / pageInfo.pageSize) + (pageInfo.count % pageInfo.pageSize > 0 ? 1 : 0)
    // console.log(pageInfo.pageNum, pageInfo.pageCount, pageInfo.count)
}

const loadCategories = async() =>{
    let res = await axios.get("/category")
    // console.log(res)
    categoryOptions.value = res.data.data.categories.map((item)=>{
      return {
        label:item.name,
        value:item.id
      }
    })
}

const categoryName = computed(() => {
    let selectedOption = categoryOptions.value.find((option) => {return option.value == selectedCategory.value})
    // console.log(selectedOption)
    return selectedOption ? selectedOption.label : ""
})

const searchByCategory = (categoryId) => {
    pageInfo.categoryId = categoryId
    pageInfo.pageNum = 1
    loadArticles()
}

const toDetail = (article) => {
    router.push({
        path: "/detail",
        query: {
            id: article.id
        }
    }) 
}

</script>

<style lang="scss" scoped>
.card {
    position: absolute;
    top: 0px;
    left: 100px;
    right: 0;
    margin: auto;
    height: 60px;
    background: white;  
    /*box-shadow: 0px 1px 3px #D3D4D8; */
    border-radius: 5px;
}
.tabs {
    position: absolute;
    top: 65px;
    left: 0;
    right: 0;
    margin: 0px auto;
    width: 1000px;
    height: 88%;
    overflow: hidden;
    overflow-y: auto;
    background: white;  
    box-shadow: 0px 1px 3px #D3D4D8; 
    border-radius: 5px; 
}

::-webkit-scrollbar {
  width:5px;
}

.card-text {
  width: 100%;
  text-align: center;
  height: 200px;
  line-height: 200px;
  font-size: 30px;
  cursor: pointer;
  color: dodgerblue;
}

.other-info-part {
  /*position: absolute;*/
  margin-top: 10px;
  display: inline-block;
  margin-right: 30px;
}
</style>