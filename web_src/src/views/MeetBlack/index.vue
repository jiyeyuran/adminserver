<template>
  <div class="container_meetIndex">
    <!-- 头部 -->
    <!-- <div class="meetIndex_header">
      <div class="btn_suc" @click="meetAdd">
        <i class="iconfont iconadd"></i>
        创建
      </div> -->
      <!-- <div class="btn_err">
        <i class="iconfont icondelete"></i>
        删除
      </div>-->
      <!-- <div class="inp_search">
        <input type="text" placeholder="请输入搜索内容" />
        <i class="iconfont iconsearch"></i>
      </div>
    </div> -->
    <!-- 主体 -->
    <div class="meetIndex_body Table">
      <el-table :data="videoList" style="width: 100%" stripe>
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="roomName" label="会议室名称" ></el-table-column>
        <el-table-column prop="ctime" label="开始时间" ></el-table-column>
        <el-table-column prop="duration" label="录制时长"></el-table-column>
        <el-table-column prop="size" label="文件大小"></el-table-column>
        <el-table-column label="操作" fixed="right" width="130">
          <template slot-scope="scope">
            <!-- <span class="suc_Col operation" @click="startMeet(scope.row)">
              <i class="iconfont iconbofang1"></i>
            </span> -->

            <!-- <router-link :to="{path:'/editMeet',query:{id:scope.row.id}}" class="suc_Col operation">
              <i class="iconfont iconxiazai1"></i>
            </router-link> -->

            <span class="del_Col operation" @click="deleVideoL(scope.row)">
              <i class="iconfont icondelete"></i>
            </span>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="meetIndex_footer pagination">
      <el-pagination
        background
        :current-page.sync="page"
        :page-size="perPage"
        layout="total, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </div>
  </div>
</template>
<script>
import { getVideoList,deleVideoL } from "../../request/modules/meetback";
export default {
  data() {
    return {
      total: 0,
      page: 1,
			perPage: 10,
			videoList:[]
    };
  },
  watch: {
    page: function () {
      this.getVideoList();
    },
  },
  mounted() {
    this.getVideoList();
  },
  methods: {
    // 获取录像列表
    getVideoList() {
      getVideoList({
        page: this.page - 1,
        perPage: this.perPage,
      }).then((res) => {
				console.log(res);
        this.total = res.count;
				this.videoList=res.items
				console.log(this.videoList);
      });
    },
    // 删除录像
    deleVideoL(e) {
      this.$confirm(`确认删除吗？`, "提示")
        .then(() => {
          deleVideoL({ id: e.id }).then((res) => {
            this.$message({
              message: "删除成功",
              type: "success",
            });
            this.getVideoList();
          });
        })
        .catch(() => {});
    },
  },
};
</script>


