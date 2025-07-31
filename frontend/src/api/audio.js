import request from '@/utils/request'

export function uploadAudio(data) {
  return request({
    url: '/audios/upload',
    method: 'post',
    data
  })
}

export function listAudios() {
  return request({
    url: '/audios',
    method: 'get'
  })
}

export function deleteAudio(id) {
  return request({
    url: `/audios/${id}`,
    method: 'delete'
  })
}