<template>
  <div>
    <a-card>
      <h3>{{ id ? '编辑文章' : '新增文章' }}</h3>

      <a-form :model="artInfo" ref="formRef" :rules="rules"
              :hideRequiredMark="true">
        <a-row :gutter="24">
          <a-col :span="16">
            <a-form-item label="文章标题" name="title">
              <a-input style="width: 300px" v-model:value="artInfo.title"></a-input>
            </a-form-item>
            <a-form-item label="文章描述" name="desc">
              <a-input type="textarea" v-model:value="artInfo.desc"></a-input>
            </a-form-item>
          </a-col>
          <a-col :span="8">

            <a-form-item label="文章分类" name="category_id">
              <article-category-select
                  style="width: 130px"
                  placeholder="请选择分类"
                  v-model:val="artInfo.category_id"
              />
            </a-form-item>

            <a-form-item label="文章缩略图" name="img">
              <file-upload
                  accept="image/jpeg,image/jpg,image/png"
                  @change="uploadChange"
              />
            </a-form-item>

          </a-col>
        </a-row>

        <a-form-item name="content">
          <Toolbar style="border-bottom: 1px solid #ccc" :editor="editorRef" :defaultConfig="toolbarConfig"
                   mode="default"/>
          <Editor
              style="height: 440px; overflow-y: hidden"
              v-model="artInfo.content"
              :defaultConfig="editorConfig"
              mode="default"
              @onCreated="handleCreated"
          ></Editor>
        </a-form-item>

        <a-form-item>
          <a-button type="danger" style="margin-right: 15px" @click="addArticle(artInfo.id)">
            {{
              artInfo.id ? '更新' : '提交'
            }}
          </a-button>
          <a-button type="primary" @click="addCancel">取消</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type {FormInstance} from 'ant-design-vue';
import {message} from "ant-design-vue";
import {ref, shallowRef} from 'vue';
import type {Rule} from 'ant-design-vue/es/form';
import {useRouter} from "vue-router";
import {getToken} from "@/plugin/utils/cache";
import "@wangeditor/editor/dist/css/style.css"
import {Editor, Toolbar} from "@wangeditor/editor-for-vue"
import ArticleCategorySelect from "../global_components/article_category_select.vue"
import FileUpload from "../global_components/file_upload.vue"
import {lbblogOnece} from "@/plugin/api/lbbill";

// 路由声明
const router = useRouter()

const artInfo = ref<lbblog.ModelArticle>({});

const rules: Record<string, Rule[]> = {
  title: [{required: true, message: '请输入文章标题', trigger: 'change'}],
  category_id: [{required: true, message: '请选择文章分类', trigger: 'change'}],
  desc: [
    {required: true, message: '请输入文章描述', trigger: 'change'},
    {max: 120, message: '描述最多可写120个字符', trigger: 'change'},
  ],
  img: [{required: true, message: '请选择文章缩略图', trigger: 'change'}],
  content: [{required: true, message: '请输入文章内容', trigger: 'change'}],
};

const id = Number(router.currentRoute.value.query.id)

const getArticle = async (id: number | undefined) => {
  if (!id) {
    return
  }
  try {
    const resp = await lbblogOnece.getArticle({
      id: id,
    })
    artInfo.value = resp.article
  } catch (error: any) {
    message.error(error);
  }
}
getArticle(id)

const uploadChange = async (data: string) => {
  artInfo.value!.img = data
}

const addArticle = async (id: number | undefined) => {
  try {
    if (id) {
      await lbblogOnece.updateArticle({
        article: artInfo.value,
      })
    } else {
      await lbblogOnece.addArticle({
        article: artInfo.value,
      })
    }
  } catch (error: any) {
    message.error(error);
  }
  // 跳转路由
  await router.push({
    name: "article_list",
  })
};

const formRef = ref<FormInstance>();
const addCancel = async () => {
  formRef.value!.resetFields();
  // 跳转路由
  await router.push({
    name: "article_list",
  })
};

const headers = {Authorization: getToken()}

const editorRef = shallowRef() // 编辑器实例，必须用 shallowRef
const valueHtml = ref("") // 内容 HTML
const toolbarConfig = {
  excludeKeys: ["insertLink", "insertImage", "editImage", "viewImageLink", "group-video", "emotion", "fullScreen"],
}
const editorConfig = {
  placeholder: "请输入内容...", MENU_CONF: {}
}
const handleCreated = (editor: any) => {
  editorRef.value = editor // 记录 editor 实例，重要！
}

</script>
