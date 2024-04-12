<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryForm" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="关键词" prop="configName">
        <el-input
          v-model="queryParams.keyword"
          placeholder="请输入文章关键词"
          clearable
          size="small"
          style="width: 240px"
          @keyup.enter.native="handleQuery"
        />
      </el-form-item>

      <el-form-item label="发布时间">
        <el-date-picker
          v-model="dateRange"
          size="small"
          style="width: 240px"
          value-format="yyyy-MM-dd"
          type="daterange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
        ></el-date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-search" size="mini" @click="handleQuery">搜索</el-button>
        <el-button icon="el-icon-refresh" size="mini" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button
          type="success"
          plain
          icon="el-icon-search"
          size="mini"
          :disabled="multiple"
          @click="handleAudit"
          v-hasPermi="['system:article:audit']"
        >审核</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="danger"
          plain
          icon="el-icon-delete"
          size="mini"
          :disabled="multiple"
          @click="handleDelete"
          v-hasPermi="['system:config:remove']"
        >删除</el-button>
      </el-col>
      <right-toolbar :showSearch.sync="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="articleList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="文章ID" align="center" prop="id" />
      <el-table-column label="文章标题" align="center" prop="title" :show-overflow-tooltip="true" />
      <el-table-column label="发布人" align="center" prop="createUser" />
      <el-table-column label="状态" align="center" width="100">
        <template slot-scope="scope">
          <span v-if="scope.row.status==0" style="color:blue;">待审核</span>
          <span v-if="scope.row.status==1" style="color:green;">正常</span>
          <span v-if="scope.row.status==2" style="color:red;">审核不通过</span>
        </template>
      </el-table-column>
      <el-table-column label="发布时间" align="center" prop="created_at" width="180">
        <template slot-scope="scope">
          <span>{{ $moment(scope.row.created_at).format('YYYY-MM-DD HH:mm:ss')  }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template slot-scope="scope">
          <el-button
            size="mini"
            type="text"
            icon="el-icon-view"
            @click="handleView(scope.row)"
          >查看</el-button>
          <el-button
            size="mini"
            type="text"
            icon="el-icon-search"
            @click="handleAudit(scope.row)"
            v-hasPermi="['system:article:audit']"
          >审核</el-button>
          <el-button
            size="mini"
            type="text"
            icon="el-icon-delete"
            @click="handleDelete(scope.row)"
            v-hasPermi="['system:article:remove']"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="queryParams.pageNum"
      :limit.sync="queryParams.pageSize"
      @pagination="getList"
    />

    <!-- 添加或修改文章对话框 -->
    <el-dialog :title="title" :visible.sync="open" width="500px" append-to-body>
      <div class="article-info">
        <div>标题：{{ articleInfo.title }}</div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleAudit">确 定</el-button>
        <el-button @click="cancle">取 消</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {getArticle, listArticle, delArticle, auditArticle} from "@/api/system/article";

export default {
  name: "Article",
  data() {
    return {
      // 遮罩层
      loading: true,
      // 选中数组
      ids: [],
      // 非单个禁用
      single: true,
      // 非多个禁用
      multiple: true,
      // 显示搜索条件
      showSearch: true,
      // 总条数
      total: 0,
      // 文章表格数据
      articleList: [],
      // 弹出层标题
      title: "",
      // 是否显示弹出层
      open: false,
      // 日期范围
      dateRange: [],
      // 查询文章
      queryParams: {
        pageNum: 1,
        pageSize: 10,
        keyword: undefined
      },
      articleInfo: {},
    };
  },
  created() {
    this.getList();
  },
  methods: {
    /** 查询文章列表 */
    getList() {
      this.loading = true;
      listArticle(this.addDateRange(this.queryParams, this.dateRange)).then(response => {
        let articles = response.data.list;
        let articleUsers = response.data.publish_users;
        for (let index in articles) {
          let article = articles[index]
          articles[index]["createUser"] = articleUsers[article.user_id].userName
        }
          this.articleList = articles;
          this.total = response.data.total;
          this.loading = false;
        }
      );
    },
    /** 搜索按钮操作 */
    handleQuery() {
      this.queryParams.pageNum = 1;
      this.getList();
    },
    /** 重置按钮操作 */
    resetQuery() {
      this.dateRange = [];
      this.resetForm("queryForm");
      this.handleQuery();
    },
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.id)
      this.single = selection.length!=1
      this.multiple = !selection.length
    },
    /** 修改按钮操作 */
    handleAudit(row) {
      const articleId = row.id || this.ids
      this.$confirm('是否确认审核通过文章编号为"' + articleId + '"的数据项?', "警告", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(function() {
        return auditArticle({'articleId': articleId});
      }).then(() => {
        this.getList();
        this.msgSuccess("审核已通过");
      }).catch(() => {});
    },
    cancle() {
      this.open = false;
    },
    handleView(row){
      const articleId = row.id || this.ids
      getArticle(articleId).then(response => {
        this.articleInfo = response.data
        this.open = true
        this.title = "查看文章详情"
      })
    },
    /** 删除按钮操作 */
    handleDelete(row) {
      const articleIds = row.id || this.ids;
      this.$confirm('是否确认删除文章编号为"' + articleIds + '"的数据项?', "警告", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }).then(function() {
          return delArticle(articleIds);
        }).then(() => {
          this.getList();
          this.msgSuccess("删除成功");
        }).catch(() => {});
    }
  }
};
</script>
