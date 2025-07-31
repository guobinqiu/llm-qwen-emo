<template>
  <div class="audio-library">
    <el-card class="main-card">
      <div slot="header" class="clearfix">
        <span>音频管理</span>
      </div>

      <!-- 上传区域 -->
      <el-upload
        class="upload-demo"
        drag
        :action="apiUrl + '/audios/upload'"
        :headers="uploadHeaders"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeAudioUpload"
        multiple
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <div class="el-upload__tip" slot="tip">
          支持 wav、mp3 格式，音频文件大小&lt;15M
        </div>
      </el-upload>

      <!-- 音频列表 -->
      <div class="audio-list" v-if="audioList.length > 0">
        <el-row :gutter="20">
          <el-col
            :span="6"
            v-for="audio in audioList"
            :key="audio.id"
            class="audio-item"
          >
            <el-card :body-style="{ padding: '10px' }">
              <div class="audio-info">
                <div class="audio-name">{{ audio.filename }}</div>
                <div class="audio-meta">
                  <div class="audio-size">{{ formatFileSize(audio.size) }}</div>
                  <el-slider
                    v-if="currentUrl === audio.url"
                    v-model="currentTime"
                    :max="duration"
                    @change="onSeek"
                    size="small"
                    style="width: 100%; margin-top: 5px"
                  />
                </div>
                <div class="audio-actions">
                  <el-button
                    type="primary"
                    :icon="isPlaying(audio.url) ? 'el-icon-video-pause' : 'el-icon-video-play'"
                    circle
                    size="mini"
                    @click="playAudio(audio.url)"
                  ></el-button>
                  <el-button
                    type="danger"
                    icon="el-icon-delete"
                    circle
                    size="mini"
                    @click="deleteAudio(audio.id)"
                  ></el-button>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <div v-else class="no-data">
        暂无音频数据
      </div>
    </el-card>
  </div>
</template>

<script>
import { listAudios, deleteAudio } from '@/api/audio'
import { getToken } from '@/utils/auth'
import { Message } from 'element-ui'

export default {
  name: 'AudioLibrary',
  data() {
    return {
      audioList: [],
      currentAudio: null,
      currentUrl: null,
      currentTime: 0,
      duration: 0,
      timer: null,
      loading: false,
    }
  },
  computed: {
    uploadHeaders() {
      return {
        Authorization: 'Bearer ' + getToken(),
      }
    },
    apiUrl() {
      return 'http://localhost:8080/api'
    },
  },
  mounted() {
    this.fetchAudioList()
  },
  methods: {
    fetchAudioList() {
      this.loading = true
      listAudios()
        .then((res) => {
          this.audioList = res.data.audios || []
        })
        .catch((err) => {
          console.error('获取音频列表失败:', err)
          Message({
            message: '获取音频列表失败',
            type: 'error',
            duration: 3000,
          })
        })
        .finally(() => {
          this.loading = false
        })
    },
    handleUploadSuccess() {
      Message({
        message: '上传成功',
        type: 'success',
        duration: 2000,
      })
      this.fetchAudioList()
    },
    handleUploadError(err) {
      console.error('上传失败:', err)
      Message({
        message: '上传失败，请检查文件格式和大小',
        type: 'error',
        duration: 3000,
      })
    },
    beforeAudioUpload(file) {
      const isAudio = ['audio/wav', 'audio/mpeg'].includes(file.type)
      const isLt15M = file.size / 1024 / 1024 < 15

      if (!isAudio) {
        Message({
          message: '上传音频只能是 WAV/MP3 格式!',
          type: 'error',
          duration: 3000,
        })
        return false
      }
      if (!isLt15M) {
        Message({
          message: '上传音频大小不能超过 15MB!',
          type: 'error',
          duration: 3000,
        })
        return false
      }
      return true
    },
    deleteAudio(id) {
      this.$confirm('此操作将永久删除该音频, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(() => {
          deleteAudio(id)
            .then(() => {
              Message({
                message: '删除成功',
                type: 'success',
                duration: 2000,
              })
              this.fetchAudioList()
            })
            .catch((err) => {
              console.error('删除失败:', err)
              Message({
                message: '删除失败',
                type: 'error',
                duration: 3000,
              })
            })
        })
        .catch(() => {
          Message({
            message: '已取消删除',
            type: 'info',
            duration: 1500,
          })
        })
    },
    playAudio(url) {
      if (this.currentAudio && this.currentUrl === url) {
        if (this.currentAudio.paused) {
          this.currentAudio.play()
        } else {
          this.currentAudio.pause()
        }
      } else {
        if (this.currentAudio) {
          this.currentAudio.pause()
          clearInterval(this.timer)
        }
        this.currentAudio = new Audio(url)
        this.currentUrl = url
        this.currentAudio.play().catch((err) => {
          console.error('播放失败:', err)
          Message({
            message: '无法播放音频',
            type: 'error',
            duration: 2000,
          })
        })
        this.currentAudio.addEventListener('loadedmetadata', () => {
          this.duration = this.currentAudio.duration
        })
        this.currentAudio.addEventListener('ended', () => {
          this.currentAudio = null
          this.currentUrl = null
          this.currentTime = 0
          this.duration = 0
          clearInterval(this.timer)
        })
        this.timer = setInterval(() => {
          if (this.currentAudio) this.currentTime = this.currentAudio.currentTime
        }, 500)
      }
    },
    onSeek(val) {
      if (this.currentAudio) {
        this.currentAudio.currentTime = val
      }
    },
    isPlaying(url) {
      return this.currentAudio && this.currentUrl === url && !this.currentAudio.paused
    },
    formatFileSize(size) {
      if (!size) return '0B'
      const units = ['B', 'KB', 'MB', 'GB']
      let unitIndex = 0
      let fileSize = size
      while (fileSize >= 1024 && unitIndex < units.length - 1) {
        fileSize /= 1024
        unitIndex++
      }
      return fileSize.toFixed(1) + units[unitIndex]
    },
  },
  beforeDestroy() {
    if (this.currentAudio) {
      this.currentAudio.pause()
      clearInterval(this.timer)
    }
  },
}
</script>

<style scoped>
.main-card {
  min-height: calc(100vh - 100px);
}
.audio-library {
  padding: 0;
}

.upload-demo {
  margin-bottom: 30px;
}

.audio-list {
  margin-top: 20px;
}

.audio-item {
  margin-bottom: 20px;
}

.audio-info {
  display: flex;
  flex-direction: column;
}

.audio-name {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 5px;
}

.audio-meta {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.audio-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
