import request from '@/utils/request'

export function createTask(data) {
  return request({
    url: '/tasks',
    method: 'post',
    data
  })
}

export function getTask(id) {
  return request({
    url: `/tasks/${id}`,
    method: 'get'
  })
}

export function listTasks() {
  return request({
    url: '/tasks',
    method: 'get'
  })
}

export function deleteTask(id) {
  return request({
    url: `/tasks/${id}`,
    method: 'delete'
  })
}