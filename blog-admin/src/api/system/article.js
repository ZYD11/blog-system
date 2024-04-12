import request from '@/utils/request'

// 查询文章列表
export function listArticle(query) {
  return request({
    url: '/api/admin/article/list',
    method: 'get',
    params: query
  })
}

// 查询文章详细
export function getArticle(articleId) {
  return request({
    url: '/api/admin/article/' + articleId,
    method: 'get'
  })
}


// 删除文章
export function delArticle(articleId) {
  return request({
    url: '/api/admin/article/' + articleId,
    method: 'delete'
  })
}

// 文章审核
export function auditArticle(data) {
  return request({
    url: '/api/admin/article/audit',
    method: 'post',
    data: data
  })
}
