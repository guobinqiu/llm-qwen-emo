import request from '@/utils/request'

export function uploadImage(data) {
  return request({
    url: '/images/upload',
    method: 'post',
    data
  })
}

export function listImages() {
  return request({
    url: '/images',
    method: 'get'
    })
}

export function deleteImage(id) {
  return request({
    url: `/images/${id}`,
    method: 'delete'
  })
}