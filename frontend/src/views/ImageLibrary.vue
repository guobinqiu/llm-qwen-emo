<template>
  <div class="image-library">
    <el-card class="main-card">
      <div slot="header" class="clearfix">
        <span>图像管理</span>
      </div>

      <!-- 上传区域 -->
      <el-upload
        class="upload-demo"
        drag
        :action="apiUrl + '/images/upload'"
        :headers="uploadHeaders"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeImageUpload"
        multiple
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <div class="el-upload__tip" slot="tip">
          支持 jpg、jpeg、png、bmp、webp 格式，图像最小边长≥400像素，最大边长≤7000像素
        </div>
      </el-upload>

      <!-- 图像列表 -->
      <div class="image-list" v-if="imageList.length > 0">
        <el-row :gutter="20">
          <el-col :span="6" v-for="image in imageList" :key="image.id" class="image-item">
            <el-card :body-style="{ padding: '10px' }">
              <img
                :src="image.url"
                class="image-preview"
                alt="预览图"
                @click="handlePreview(image)"
                style="cursor: pointer;"
              />
              <div class="image-info">
                <div class="image-name">{{ image.name }}</div>
                <el-button
                  type="danger"
                  icon="el-icon-delete"
                  circle
                  size="mini"
                  @click="deleteImage(image.id)"
                ></el-button>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <div v-else class="no-data">
        暂无图像数据
      </div>
    </el-card>

    <!-- 预览弹窗 -->
    <el-dialog
      :visible.sync="previewVisible"
      width="60%"
      center
      :show-close="true"
      custom-class="image-preview-dialog"
    >
      <img :src="previewImageUrl" class="dialog-image" alt="大图预览" />
    </el-dialog>
  </div>
</template>

<script>
import { listImages, deleteImage } from '@/api/image'
import { getToken } from '@/utils/auth'
import { Message } from 'element-ui'

export default {
  name: 'ImageLibrary',
  data() {
    return {
      imageList: [],
      previewVisible: false,
      previewImageUrl: ''
    }
  },
  computed: {
    uploadHeaders() {
      return {
        Authorization: 'Bearer ' + getToken()
      }
    },
    apiUrl() {
      return 'http://localhost:8080/api'
    }
  },
  mounted() {
    this.fetchImageList()
  },
  methods: {
    fetchImageList() {
      listImages()
        .then((res) => {
          this.imageList = res.data.images || []
        })
        .catch((err) => {
          console.error('获取图像列表失败:', err)
          Message({
            message: '获取图像列表失败',
            type: 'error',
            duration: 3000
          })
        })
    },
    handleUploadSuccess() {
      Message({
        message: '上传成功',
        type: 'success',
        duration: 2000
      })
      this.fetchImageList()
    },
    handleUploadError(err) {
      console.error('上传失败:', err)
      Message({
        message: '上传失败，请检查文件格式和大小',
        type: 'error',
        duration: 3000
      })
    },
    beforeImageUpload(file) {
      const isImage = ['image/jpeg', 'image/png', 'image/bmp', 'image/webp'].includes(file.type)
      const isLt10M = file.size / 1024 / 1024 < 10

      if (!isImage) {
        Message({
          message: '上传图像只能是 JPG/PNG/BMP/WEBP 格式!',
          type: 'error',
          duration: 3000
        })
        return false
      }
      if (!isLt10M) {
        Message({
          message: '上传图像大小不能超过 10MB!',
          type: 'error',
          duration: 3000
        })
        return false
      }
      return true
    },
    deleteImage(id) {
      this.$confirm('此操作将永久删除该图像, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          deleteImage(id)
            .then(() => {
              Message({
                message: '删除成功',
                type: 'success',
                duration: 2000
              })
              this.fetchImageList()
            })
            .catch((err) => {
              console.error('删除失败:', err)
              Message({
                message: '删除失败',
                type: 'error',
                duration: 3000
              })
            })
        })
        .catch(() => {
          Message({
            message: '已取消删除',
            type: 'info',
            duration: 1500
          })
        })
    },
    handlePreview(image) {
      this.previewImageUrl = image.url
      this.previewVisible = true
    }
  }
}
</script>

<style scoped>
.image-library {
  padding: 0;
}

.main-card {
  min-height: calc(100vh - 100px);
}

.upload-demo {
  margin-bottom: 30px;
}

.image-list {
  margin-top: 20px;
}

.image-item {
  margin-bottom: 20px;
}

.image-preview {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.image-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.image-name {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

.image-preview-dialog >>> .el-dialog__body {
  padding: 0;
  text-align: center;
}

.dialog-image {
  max-width: 100%;
  max-height: 80vh;
  display: block;
  margin: 0 auto;
}
</style>
