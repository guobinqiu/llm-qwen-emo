<template>
  <div class="task-management">
    <el-card>
      <div slot="header" class="clearfix">
        <span>任务管理</span>
      </div>

      <!-- 创建任务表单 -->
      <el-form :model="taskForm" :rules="taskRules" ref="taskForm" label-width="100px" class="task-form">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="选择图像" prop="imageId">
          <el-select
            v-model="taskForm.imageId"
            placeholder="请选择图像"
            class="wide-select"
            popper-class="wide-option-popper"
          >
            <el-option
              v-for="image in imageList"
              :key="image.id"
              :label="image.filename"
              :value="image.id"
            >
              <div class="option-item">
                <span class="left">{{ image.filename }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="选择音频" prop="audioId">
          <el-select
            v-model="taskForm.audioId"
            placeholder="请选择音频"
            class="wide-select"
            popper-class="wide-option-popper"
          >
            <el-option
              v-for="audio in audioList"
              :key="audio.id"
              :label="audio.filename"
              :value="audio.id"
            >
              <div class="option-item">
                <span class="left">{{ audio.filename }}</span>
                <span class="right">{{ formatFileSize(audio.size) }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="分段时长" prop="segmentSecond">
          <el-input-number
            v-model="taskForm.segmentSecond"
            :min="1"
            :max="59"
            label="分段时长"
            style="width: 150px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm('taskForm')">创建任务</el-button>
          <el-button @click="resetForm('taskForm')">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 任务列表 -->
      <div class="task-list">
        <h3>任务列表</h3>
        <el-table
          :data="taskList"
          style="width: 100%"
          v-loading="loading"
          row-key="id"
        >
          <el-table-column type="expand">
            <template #default="{ row }">
              <el-table
                :data="row.sub_tasks || []"
                size="small"
                style="width: 100%; margin-left: 20px"
                row-key="code"
              >
                <el-table-column prop="task_status" label="子任务状态" />
                <el-table-column prop="scheduled_time" label="调度时间">
                  <template #default="{ row }">
                    {{ formatDate(row.scheduled_time) }}
                  </template>
                </el-table-column>
                <el-table-column prop="end_time" label="结束时间">
                  <template #default="{ row }">
                    {{ formatDate(row.end_time) }}
                  </template>
                </el-table-column>
                <el-table-column prop="message" label="错误信息" />
                <el-table-column prop="audio_url" label="音频URL" />
                <el-table-column prop="video_url" label="视频URL" />
                <el-table-column label="操作">
                  <template #default="{ row }">
                    <el-button
                      size="mini"
                      type="success"
                      :disabled="!row.video_url"
                      @click="playVideo(row.video_url)"
                    >
                      播放视频
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </template>
          </el-table-column>

          <el-table-column prop="name" label="任务名称" />
          <el-table-column prop="status" label="任务状态" />
          <el-table-column prop="message" label="错误信息" />
          <el-table-column prop="created_at" label="创建时间">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="image_url" label="图像URL" />
          <el-table-column prop="video_url" label="视频URL" />

          <el-table-column label="操作">
            <template #default="{ row }">
              <el-button
                size="mini"
                type="success"
                :disabled="!row.video_url"
                @click="playVideo(row.video_url)"
              >
                播放视频
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 视频播放弹窗 -->
    <el-dialog :visible.sync="videoDialogVisible" title="视频预览" width="60%">
      <video v-if="videoUrl" controls style="width: 100%">
        <source :src="videoUrl" type="video/mp4" />
        您的浏览器不支持 video 标签。
      </video>
    </el-dialog>
  </div>
</template>

<script>
import { listImages } from '@/api/image'
import { listAudios } from '@/api/audio'
import { createTask, listTasks } from '@/api/task'
import { Message } from 'element-ui'

export default {
  name: 'TaskManagement',
  data() {
    return {
      taskForm: {
        name: '',
        imageId: '',
        audioId: '',
        segmentSecond: 15,
        timer: null,
      },
      taskRules: {
        name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
        imageId: [{ required: true, message: '请选择图像', trigger: 'change' }],
        audioId: [{ required: true, message: '请选择音频', trigger: 'change' }],
        segmentSecond: [
          { required: true, message: '请输入分段时长', trigger: 'change' },
          {
            type: 'number',
            min: 1,
            max: 59,
            message: '分段时长必须在 1 到 59 秒之间',
            trigger: 'change'
          }
        ]
      },
      imageList: [],
      audioList: [],
      taskList: [],
      loading: false,
      videoDialogVisible: false,
      videoUrl: ''
    }
  },
  mounted() {
    this.fetchImageList()
    this.fetchAudioList()
    this.fetchTaskList()
    this.timer = setInterval(() => {
      this.fetchTaskList()
    }, 10000)
  },
  beforeDestroy() {
    if (this.timer) {
      clearInterval(this.timer)
    }
  },
  methods: {
    fetchImageList() {
      listImages().then(res => {
        this.imageList = res.data.images || []
      })
    },
    fetchAudioList() {
      listAudios().then(res => {
        this.audioList = res.data.audios || []
      })
    },
    fetchTaskList() {
      this.loading = true
      listTasks()
        .then(res => {
          this.taskList = res.data.tasks || []
          this.loading = false
        })
        .catch(() => {
          this.loading = false
          Message.error('获取任务列表失败')
        })
    },
    submitForm(formName) {
      this.$refs[formName].validate(valid => {
        if (!valid) return
        const data = {
          name: this.taskForm.name,
          image_id: this.taskForm.imageId,
          audio_id: this.taskForm.audioId,
          segment_second: this.taskForm.segmentSecond
        }
        createTask(data)
          .then(() => {
            Message.success('任务创建成功')
            this.resetForm(formName)
            this.fetchTaskList()
          })
          .catch(() => {
            Message.error('创建任务失败')
          })
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    playVideo(url) {
      this.videoUrl = url
      this.videoDialogVisible = true
    },
    formatDate(dateVal) {
      if (!dateVal) return ''
      const date = new Date(dateVal)
      if (isNaN(date)) return ''
      return date.toLocaleString('zh-CN')
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
    }
  }
}
</script>

<style scoped>
.task-management {
  padding: 0;
}
.task-form {
  margin-bottom: 30px;
}
.task-list {
  margin-top: 30px;
}
.wide-select {
  width: 100%;
}
::v-deep .wide-option-popper {
  min-width: 500px !important;
}
.option-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.option-item .left {
  font-weight: 500;
  font-size: 14px;
  word-break: break-word;
}
.option-item .right {
  color: #999;
  font-size: 13px;
  margin-left: 20px;
  flex-shrink: 0;
}
</style>
