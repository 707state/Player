<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { musicApi } from "../services/api";
import type { Album, DialogFormData } from "../services/types";

const albums = ref<Album[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const isEditing = ref(false);
const currentAlbum = ref<Album | null>(null);

const form = reactive<DialogFormData>({
    title: "",
    artist: "",
    genre: "",
    year: new Date().getFullYear(),
    url: "",
    artwork: "",
    comment: "",
    rating: 0,
    cuts: [],
});

const newCut = ref("");

// Load albums
const loadAlbums = async () => {
    loading.value = true;
    try {
        albums.value = await musicApi.getAlbums();
    } catch (error) {
        ElMessage.error("加载专辑失败");
        console.error(error);
    } finally {
        loading.value = false;
    }
};

// Open add dialog
const openAddDialog = () => {
    isEditing.value = false;
    currentAlbum.value = null;
    resetForm();
    dialogVisible.value = true;
};

// Open edit dialog
const openEditDialog = (album: Album) => {
    isEditing.value = true;
    currentAlbum.value = album;
    Object.assign(form, {
        title: album.title,
        artist: album.artist,
        genre: album.genre,
        year: album.year,
        url: album.url,
        artwork: album.artwork,
        comment: album.comment,
        rating: album.rating,
        cuts: [...album.cuts],
    });
    dialogVisible.value = true;
};

// Reset form
const resetForm = () => {
    Object.assign(form, {
        title: "",
        artist: "",
        genre: "",
        year: new Date().getFullYear(),
        url: "",
        artwork: "",
        comment: "",
        rating: 0,
        cuts: [],
    });
    newCut.value = "";
};

// Add cut
const addCut = () => {
    if (newCut.value.trim()) {
        form.cuts!.push(newCut.value.trim());
        newCut.value = "";
    }
};

// Remove cut
const removeCut = (index: number) => {
    form.cuts!.splice(index, 1);
};

// Save album
const saveAlbum = async () => {
    if (!form.title || !form.artist) {
        ElMessage.warning("请填写标题和艺术家");
        return;
    }

    try {
        const albumData: Album = {
            title: form.title,
            artist: form.artist!,
            genre: form.genre,
            year: form.year,
            url: form.url,
            artwork: form.artwork || "",
            comment: form.comment,
            rating: form.rating,
            cuts: form.cuts || [],
        };

        await musicApi.saveAlbum(albumData);
        ElMessage.success(isEditing.value ? "更新成功" : "添加成功");
        dialogVisible.value = false;
        await loadAlbums();
    } catch (error) {
        ElMessage.error("保存失败");
        console.error(error);
    }
};

// Delete album
const deleteAlbum = async (album: Album) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除专辑 "${album.title}" 吗？`,
            "确认删除",
            {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning",
            },
        );

        await musicApi.deleteAlbum({
            title: album.title,
            artist: album.artist,
        });
        ElMessage.success("删除成功");
        await loadAlbums();
    } catch (error) {
        if (error !== "cancel") {
            ElMessage.error("删除失败");
            console.error(error);
        }
    }
};

// Initialize
onMounted(() => {
    loadAlbums();
});
</script>

<template>
    <div class="music-page">
        <div class="page-header">
            <h2>音乐收藏</h2>
            <el-button type="primary" @click="openAddDialog" class="add-button">
                <el-icon><Plus /></el-icon>
                添加专辑
            </el-button>
        </div>

        <div class="albums-grid" v-loading="loading">
            <div
                v-for="album in albums"
                :key="`${album.title}-${album.artist}`"
                class="album-card"
            >
                <div class="album-image">
                    <img
                        v-if="album.artwork"
                        :src="album.artwork"
                        :alt="album.title"
                        @error="(e: any) => (e.target.style.display = 'none')"
                    />
                    <div v-else class="no-image">
                        <el-icon><Picture /></el-icon>
                    </div>
                </div>
                <div class="album-info">
                    <h3 class="album-title">{{ album.title }}</h3>
                    <p class="album-artist">{{ album.artist }}</p>
                    <div class="album-meta">
                        <span class="genre">{{ album.genre }}</span>
                        <span class="year">{{ album.year }}</span>
                    </div>
                    <div class="album-rating">
                        <el-rate
                            v-model="album.rating"
                            disabled
                            show-score
                            text-color="#ff9900"
                            score-template="{value}"
                        />
                    </div>
                    <p v-if="album.comment" class="album-comment">
                        {{ album.comment }}
                    </p>
                    <div
                        v-if="album.cuts && album.cuts.length > 0"
                        class="album-cuts"
                    >
                        <span class="cuts-label">曲目:</span>
                        <span class="cuts-list">{{
                            album.cuts.join(", ")
                        }}</span>
                    </div>
                </div>
                <div class="album-actions">
                    <el-button
                        type="primary"
                        size="small"
                        @click="openEditDialog(album)"
                        class="edit-btn"
                    >
                        <el-icon><Edit /></el-icon>
                        编辑
                    </el-button>
                    <el-button
                        type="danger"
                        size="small"
                        @click="deleteAlbum(album)"
                        class="delete-btn"
                    >
                        <el-icon><Delete /></el-icon>
                        删除
                    </el-button>
                </div>
            </div>

            <div v-if="albums.length === 0 && !loading" class="empty-state">
                <el-empty description="暂无专辑，点击添加按钮开始记录" />
            </div>
        </div>

        <!-- Add/Edit Dialog -->
        <el-dialog
            v-model="dialogVisible"
            :title="isEditing ? '编辑专辑' : '添加专辑'"
            width="600px"
            :before-close="() => (dialogVisible = false)"
        >
            <el-form :model="form" label-width="80px">
                <el-form-item label="标题" required>
                    <el-input
                        v-model="form.title"
                        placeholder="请输入专辑标题"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="艺术家" required>
                    <el-input
                        v-model="form.artist"
                        placeholder="请输入艺术家"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="流派">
                    <el-input v-model="form.genre" placeholder="请输入流派" />
                </el-form-item>
                <el-form-item label="年份">
                    <el-input-number
                        v-model="form.year"
                        :min="1900"
                        :max="new Date().getFullYear()"
                    />
                </el-form-item>
                <el-form-item label="评分">
                    <el-rate v-model="form.rating" show-text />
                </el-form-item>
                <el-form-item label="链接">
                    <el-input v-model="form.url" placeholder="请输入专辑链接" />
                </el-form-item>
                <el-form-item label="封面">
                    <el-input
                        v-model="form.artwork"
                        placeholder="请输入封面图片链接"
                    />
                </el-form-item>
                <el-form-item label="曲目">
                    <div class="cuts-input">
                        <el-input
                            v-model="newCut"
                            placeholder="输入曲目名称"
                            @keyup.enter="addCut"
                        >
                            <template #append>
                                <el-button @click="addCut">添加</el-button>
                            </template>
                        </el-input>
                    </div>
                    <div
                        v-if="form.cuts && form.cuts.length > 0"
                        class="cuts-list"
                    >
                        <el-tag
                            v-for="(cut, index) in form.cuts"
                            :key="index"
                            closable
                            @close="removeCut(index)"
                            class="cut-tag"
                        >
                            {{ cut }}
                        </el-tag>
                    </div>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input
                        v-model="form.comment"
                        type="textarea"
                        :rows="3"
                        placeholder="请输入备注"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" @click="saveAlbum">
                    {{ isEditing ? "更新" : "添加" }}
                </el-button>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.music-page {
    padding: 24px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e8e8e8;
}

.page-header h2 {
    margin: 0;
    color: #333;
    font-size: 1.5rem;
}

.add-button {
    display: flex;
    align-items: center;
    gap: 6px;
}

.albums-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
}

.album-card {
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    overflow: hidden;
    background: white;
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
}

.album-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.album-image {
    height: 180px;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
}

.album-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.no-image {
    color: #ccc;
    font-size: 3rem;
}

.album-info {
    padding: 16px;
    flex: 1;
}

.album-title {
    margin: 0 0 8px 0;
    font-size: 1.2rem;
    font-weight: 600;
    color: #333;
    line-height: 1.4;
}

.album-artist {
    margin: 0 0 12px 0;
    color: #666;
    font-size: 0.9rem;
}

.album-meta {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
    font-size: 0.85rem;
}

.genre {
    color: #409eff;
    background: #ecf5ff;
    padding: 2px 8px;
    border-radius: 4px;
}

.year {
    color: #909399;
}

.album-rating {
    margin-bottom: 12px;
}

.album-comment {
    margin: 8px 0;
    color: #666;
    font-size: 0.9rem;
    line-height: 1.4;
}

.album-cuts {
    margin-top: 8px;
    font-size: 0.85rem;
}

.cuts-label {
    color: #909399;
    margin-right: 4px;
}

.cuts-list {
    color: #666;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.album-actions {
    padding: 12px 16px;
    border-top: 1px solid #f0f0f0;
    display: flex;
    gap: 8px;
}

.edit-btn,
.delete-btn {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
}

.empty-state {
    grid-column: 1 / -1;
    padding: 60px 20px;
}

.cuts-input {
    margin-bottom: 12px;
}

.cuts-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.cut-tag {
    margin-bottom: 4px;
}

@media (max-width: 768px) {
    .music-page {
        padding: 16px;
    }

    .page-header {
        flex-direction: column;
        gap: 16px;
        align-items: stretch;
    }

    .albums-grid {
        grid-template-columns: 1fr;
    }

    .album-actions {
        flex-direction: column;
    }
}
</style>
