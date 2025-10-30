<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { WarningFilled } from "@element-plus/icons-vue";
const API_BASE_URL = import.meta.env.VITE_VUE_APP_BACKEND_URL;
// 类型定义，严格与后端保持一致
interface Album {
  title: string;
  artists: string[];
  genre?: string;
  year?: number;
  url?: string;
  artwork?: string;
  comment?: string;
  rating?: number;
}
// interface AlbumSingle {
//   title: string;
//   artists: string[];
//   album: string;
// }

const albums = ref<Album[]>([]);
const loading = ref(false);

const dialogVisible = ref(false);
const editingType = ref<'add' | 'edit'>('add');
const form = reactive<Album>({
  title: "",
  artists: [],
  genre: "",
  year: new Date().getFullYear(),
  url: "",
  artwork: "",
  comment: "",
  rating: 0,
});
const singles = ref<string[]>([]);
const singleInput = ref("");

const deleteDialogVisible = ref(false);
const currentDeleteAlbum = ref<Album | null>(null);

function showDeleteDialog(album: Album) {
  currentDeleteAlbum.value = album;
  deleteDialogVisible.value = true;
}
function closeDeleteDialog() {
  deleteDialogVisible.value = false;
  currentDeleteAlbum.value = null;
}
async function confirmDeleteAlbum() {
  if (!currentDeleteAlbum.value) return;
  try {
    await fetch(`${API_BASE_URL}/music`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: currentDeleteAlbum.value.title, artists: currentDeleteAlbum.value.artists })
    });
    const cuts = await getAllSingles(currentDeleteAlbum.value.title, currentDeleteAlbum.value.artists);
    for (const cut of cuts) {
      await removeSingle(cut, true, currentDeleteAlbum.value);
    }
    ElMessage.success("已删除");
    await fetchAlbums();
  } catch {
    ElMessage.error("删除失败");
  }
  closeDeleteDialog();
}

function artistsDisplay(arr?: string[]) {
  return (arr || []).join(", ");
}

// 获取所有专辑
async function fetchAlbums() {
  loading.value = true;
  try {
    const r = await fetch(`${API_BASE_URL}/music`);
    albums.value = await r.json();
  } catch {
    ElMessage.error('加载专辑失败');
  } finally {
    loading.value = false;
  }
}

// 获取某专辑的曲目（单曲）
async function fetchSingles(album: string, artists: string[]) {
  if (!album || !artists.length) { singles.value = []; return; }
  const url = new URL(`${API_BASE_URL}/single`, window.location.origin);
  url.searchParams.set('album', album);
  url.searchParams.set('artists', artists.join(','));
  const r = await fetch(url.toString(), { method: 'GET' });
  if (r.ok) {
    try {
      // 假设列表是 {singles:[AlbumSingle]}，否则兼容直接数组
      const d = await r.json();
      if (Array.isArray(d)) {
        singles.value = d.map((s: any) => typeof s === 'string' ? s : (s.title || ''));
      } else if (d.singles && Array.isArray(d.singles)) {
        singles.value = d.singles.map((s: any) => s.title);
      } else {
        singles.value = [];
      }
    } catch { singles.value = [] }
  } else {
    singles.value = [];
  }
}

// 弹窗入口
async function openDialog(album?: Album) {
  if (album) {
    editingType.value = 'edit';
    Object.assign(form, album);
    await fetchSingles(album.title, album.artists);
  } else {
    editingType.value = 'add';
    Object.assign(form, {
      title: "",
      artists: [],
      genre: "",
      year: new Date().getFullYear(),
      url: "",
      artwork: "",
      comment: "",
      rating: 0,
    });
    singles.value = [];
  }
  dialogVisible.value = true;
}

function handleDialogClose() {
  dialogVisible.value = false;
  singles.value = [];
}

async function saveAlbum() {
  if (!form.title || !form.artists.length) {
    ElMessage.warning("必须输入标题和艺人");
    return;
  }
  try {
    const res = await fetch(`${API_BASE_URL}/music`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(form),
    });
    if (!res.ok) throw new Error();
    ElMessage.success('保存成功');
    await fetchAlbums();
    dialogVisible.value = false;
  } catch { ElMessage.error('保存失败'); }
}

// 获取所有曲目
async function getAllSingles(album: string, artists: string[]): Promise<string[]> {
  const url = new URL(`${API_BASE_URL}/single`, window.location.origin);
  url.searchParams.set('album', album);
  url.searchParams.set('artists', artists.join(','));
  const r = await fetch(url.toString());
  if (!r.ok) return [];
  try {
    const d = await r.json();
    if (Array.isArray(d)) {
      return d.map((s: any) => typeof s === 'string' ? s : (s.title || ''));
    } else if (d.singles && Array.isArray(d.singles)) {
      return d.singles.map((s: any) => s.title);
    }
  } catch {}
  return [];
}

// 曲目添加/删除
async function addSingle() {
  const title = singleInput.value.trim();
  if (!title || !form.title || !form.artists.length) return;
  if (singles.value.includes(title)) {
    ElMessage.warning('已存在');
    return;
  }
  try {
    await fetch(`${API_BASE_URL}/single`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, artists: form.artists, album: form.title })
    });
    singles.value.push(title);
    singleInput.value = "";
    ElMessage.success("已添加");
  } catch { ElMessage.error("添加失败"); }
}

async function removeSingle(title: string, silent=false, ctxAlbum?: Album) {
  const album = ctxAlbum || form;
  try {
    await fetch(`${API_BASE_URL}/single`, {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, artists: album.artists, album: album.title })
    });
    if (!ctxAlbum) singles.value = singles.value.filter(s => s !== title);
    if (!silent) ElMessage.success("已删除");
  } catch { if (!silent) ElMessage.error("删除失败"); }
}

onMounted(fetchAlbums);
</script>

<template>
  <div class="music-page">
    <div class="page-header">
      <h2>音乐专辑</h2>
      <el-button type="primary" @click="openDialog()">新增专辑</el-button>
    </div>
    <div class="albums-grid" v-loading="loading">
      <div v-for="album in albums" :key="album.title+'-'+artistsDisplay(album.artists)" class="album-card">
        <div class="album-image" @click="openDialog(album)">
          <img v-if="album.artwork" :src="album.artwork" :alt="album.title" @error="e=> (e.target as HTMLImageElement).style.display='none'" />
          <div v-else class="no-image"><el-icon><Picture /></el-icon></div>
        </div>
        <div class="album-info">
          <h3 class="album-title">{{ album.title }}</h3>
          <div>{{ artistsDisplay(album.artists) }}</div>
          <span v-if="album.year || album.genre" class="album-meta">
            <span v-if="album.genre">{{ album.genre }}</span>
            <span v-if="album.year">{{ album.year }}</span>
          </span>
          <span v-if="album.rating" class="album-rating">
            <el-rate v-model="album.rating" disabled show-score text-color="#ff9900" score-template="{value}" />
          </span>
        </div>
        <div class="album-actions">
          <el-button type="danger" size="small" @click.stop="showDeleteDialog(album)">删除</el-button>
          <el-button type="primary" size="small" @click.stop="openDialog(album)">编辑/详情</el-button>
        </div>
      </div>
      <div v-if="albums.length===0&&!loading" class="empty-state">
        <el-empty description="暂无专辑" />
      </div>
    </div>
    <el-dialog v-model="dialogVisible" width="560px" :title="editingType==='add'?'新增专辑':'编辑专辑'" @close="handleDialogClose">
      <el-form label-width="85px">
        <el-form-item label="标题"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="艺人">
          <el-select
            v-model="form.artists"
            multiple filterable allow-create default-first-option
            placeholder="输入每位艺人，按回车确认"
            class="artists-select"/>
        </el-form-item>
        <el-form-item label="年份"><el-input-number v-model="form.year" :min="1900" :max="new Date().getFullYear()" /></el-form-item>
        <el-form-item label="流派"><el-input v-model="form.genre" /></el-form-item>
        <el-form-item label="封面"><el-input v-model="form.artwork" /></el-form-item>
        <el-form-item label="链接"><el-input v-model="form.url" /></el-form-item>
        <el-form-item label="评分"><el-rate v-model="form.rating" show-text /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.comment" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="曲目">
          <div style="display:flex;flex-wrap:wrap;gap:6px;width:100%">
            <el-tag v-for="cut in singles" :key="cut" closable @close="removeSingle(cut)">{{ cut }}</el-tag>
            <el-input v-model="singleInput" @keyup.enter.native="addSingle" style="width:120px" placeholder="曲名" />
            <el-button type="success" size="small" @click="addSingle">添加</el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleDialogClose">取消</el-button>
          <el-button type="primary" @click="saveAlbum">保存</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
    <el-dialog
      v-model="deleteDialogVisible"
      width="360px"
      :title="`删除专辑【${currentDeleteAlbum?.title}】`"
      @close="closeDeleteDialog"
    >
      <div style="text-align: center; color: #f56c6c; margin-bottom: 20px;">
        <el-icon><WarningFilled /></el-icon>
        <p>确定要删除专辑【{{ currentDeleteAlbum?.title }}】？</p>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDeleteDialog">取消</el-button>
        <el-button type="danger" @click="confirmDeleteAlbum">确定删除</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<style scoped>
.music-page { padding:24px; }
.page-header{
    display:flex;
    justify-content:space-between;
    align-items:center;
    margin-bottom:16px;
    padding-bottom:16px;
    border-bottom:1px solid #e8e8e8;
}
.page-header h2 {
    margin: 0;
    color: #333;
    font-size: 1.5rem;
}
.header-row{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;}
.albums-grid { display: grid; grid-template-columns: repeat(auto-fill,minmax(320px,1fr)); gap: 18px; }
.album-card { border:1px solid #e8e8e8; border-radius:8px; background:#fff; box-shadow:0 2px 8px #f4f4f4; display:flex;flex-direction:column;}
.album-image {height:180px;background:#f5f5f5;display:flex;align-items:center;justify-content:center;cursor:pointer;}
.album-image img{width:100%;height:100%;object-fit:cover;}
.no-image{color:#aaa;font-size:2rem}
.album-info{padding:14px 15px 10px 15px;flex:1;}
.album-title{margin:0 0 4px 0;font-weight:600;}
.album-meta{color:#666;}
.album-rating{margin:5px 0;}
.album-actions{display:flex;gap:10px;padding:8px 14px 12px 14px;border-top:1px solid #f0f0f0;justify-content:flex-end;}
.empty-state{grid-column:1/-1;padding:40px 10px}
.artists-select{width:100%}
@media (max-width: 768px) { .music-page{padding:6px;} .albums-grid{grid-template-columns:1fr;} .header-row{flex-direction:column;align-items:flex-start;gap:8px;} }
</style>