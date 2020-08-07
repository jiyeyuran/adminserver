<template>
  <div class="container_player" id="meet">
    <!-- <div class="player_footer">
      <span
        v-for="(item,key) in options"
        :key="key"
        :class="['footer_option',{'hangUp':item.id=='hangUp'}]"
      >
        <i :class="['iconfont',item.icon,]"></i>
        <span>{{item.lable}}</span>
      </span>
    </div>-->
  </div>
</template>

<script>
import { getMeet } from "../../request/modules/meet";
export default {
  data() {
    return {
      JitsiMeetExternal: "", //简会实例
      meetInfo: "", //房间信息
      camera: false,
      options: [
        { lable: "关闭摄像头", id: "camera", icon: "iconguanbishexiangtou1" },
        { lable: "静音", id: "mute", icon: "iconguanbiyinpin" },
        { lable: "共享屏幕", id: "sharingScreen", icon: "icongongxiangpingmu" },
        { lable: "电视墙", id: "videoWall", icon: "icondianshiqiang1" },
        { lable: "安全选项", id: "Security", icon: "iconanquan" },
        { lable: "背景模糊", id: "dim", icon: "iconmohu" },
        { lable: "画面质量", id: "frames", icon: "icongaoqing" },
        { lable: "录像", id: "video", icon: "iconluxiang1" },
        { lable: "视屏直播", id: "live", icon: "iconleft-shipinzhibo1" },
        { lable: "挂断", id: "hangUp", icon: "iconleft-shipinzhibo1" },
        { lable: "分享", id: "share", icon: "iconfenxiang" },
        { lable: "发言统计", id: "share", icon: "iconshujutongji" },
        { lable: "全体禁音", id: "allMute ", icon: "iconquantijingyin1" },
        { lable: "举手发言", id: "hands ", icon: "iconjushou1" },
        { lable: "文档分享", id: "hands ", icon: "icondocument" },
        { lable: "设置", id: "set ", icon: "iconservicemanagement" },
        { lable: "视图布局", id: "layout ", icon: "iconyemianbuju" },
        { lable: "参会人员", id: "meetUser ", icon: "iconcanhuirenyuan" },
        { lable: "消息", id: "message ", icon: "iconxinxi" },
        { lable: "全屏", id: "full ", icon: "iconkuozhan" },
      ],
    };
  },
  mounted() {
    if (this.$route.query.id) {
      this.getMeet();
    } else {
      this.$router.push("/MeetIndex");
    }
  },
  methods: {
    // 实例化简会系统
    setJME() {
      let that = this;
      this.JitsiMeetExternal = new JitsiMeetExternalAPI("room.jhmeeting.com", {
        roomName: "xiaohui",
        parentNode: document.getElementById("meet"),
        jwt:
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIqIiwianRpIjoiZWQyOTkxYmMtZTgwOC00ZTZlLWIxYTYtMjcyNGQ0NGQ1MDRiIiwiaWF0IjoxNTk2Njk0ODYwLCJpc3MiOiJqbWVldCIsInJvb20iOiJ4aWFvaHVpIn0.L80c5XQxzNBdYBPoWIVLkslXAe-nBOxbsNmCQtJkm5c",
      });
      this.JitsiMeetExternal.on("videoConferenceJoined", function () {
        console.log("videoConferenceJoined");

        this.JitsiMeetExternal.on("participantJoined", function (p) {
          console.log("participantJoined", p);
        });

        this.JitsiMeetExternal.on("participantJoined", function (p) {
          console.log("participantJoined", p);
        });
      });
      // 监听挂断事假
      this.JitsiMeetExternal.on("videoConferenceLeft", function (p) {
        that.$router.push("/MeetIndex");
      });
    },
    // 获取房间信息
    getMeet() {
      getMeet({
        id: parseInt(this.$route.query.id),
      }).then((res) => {
        this.meetInfo = res;
        this.setJME();
      });
    },
  },
};
</script>