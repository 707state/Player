<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { moviesApi } from "../services/api";
import type { Movie, DialogFormData } from "../services/types";

const movies = ref<Movie[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const isEditing = ref(false);
const currentMovie = ref<Movie | null>(null);

const form = reactive<DialogFormData>({
    title: "",
    director: "",
    genre: "",
    year: new Date().getFullYear(),
    url: "",
    comment: "",
    rating: 0,
});

// Load movies
const loadMovies = async () => {
    loading.value = true;
    try {
        movies.value = await moviesApi.getMovies();
    } catch (error) {
        ElMessage.error("加载影视失败");
        console.error(error);
    } finally {
        loading.value = false;
    }
};

// Open add dialog
const openAddDialog = () => {
    isEditing.value = false;
    currentMovie.value = null;
    resetForm();
    dialogVisible.value = true;
};

// Open edit dialog
const openEditDialog = (movie: Movie) => {
    isEditing.value = true;
    currentMovie.value = movie;
    Object.assign(form, {
        title: movie.title,
        director: movie.director,
        genre: movie.genre,
        year: movie.year,
        url: movie.url,
        comment: movie.comment,
        rating: movie.rating,
    });
    dialogVisible.value = true;
};

// Reset form
const resetForm = () => {
    Object.assign(form, {
        title: "",
        director: "",
        genre: "",
        year: new Date().getFullYear(),
        url: "",
        comment: "",
        rating: 0,
    });
};

// Save movie
const saveMovie = async () => {
    if (!form.title || !form.director) {
        ElMessage.warning("请填写标题和导演");
        return;
    }

    try {
        const movieData: Movie = {
            title: form.title,
            director: form.director!,
            genre: form.genre,
            year: form.year,
            url: form.url,
            comment: form.comment,
            rating: form.rating,
        };

        await moviesApi.saveMovie(movieData);
        ElMessage.success(isEditing.value ? "更新成功" : "添加成功");
        dialogVisible.value = false;
        await loadMovies();
    } catch (error) {
        ElMessage.error("保存失败");
        console.error(error);
    }
};

// Delete movie
const deleteMovie = async (movie: Movie) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除影视 "${movie.title}" 吗？`,
            "确认删除",
            {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning",
            },
        );

        await moviesApi.deleteMovie({
            title: movie.title,
            director: movie.director,
        });
        ElMessage.success("删除成功");
        await loadMovies();
    } catch (error) {
        if (error !== "cancel") {
            ElMessage.error("删除失败");
            console.error(error);
        }
    }
};

// Initialize
onMounted(() => {
    loadMovies();
});
</script>

<template>
    <div class="movies-page">
        <div class="page-header">
            <h2>影视收藏</h2>
            <el-button type="primary" @click="openAddDialog" class="add-button">
                <el-icon><Plus /></el-icon>
                添加影视
            </el-button>
        </div>

        <div class="movies-grid" v-loading="loading">
            <div
                v-for="movie in movies"
                :key="`${movie.title}-${movie.director}`"
                class="movie-card"
            >
                <div class="movie-image">
                    <div class="no-image">
                        <el-icon><VideoCamera /></el-icon>
                    </div>
                </div>
                <div class="movie-info">
                    <h3 class="movie-title">{{ movie.title }}</h3>
                    <p class="movie-director">{{ movie.director }}</p>
                    <div class="movie-meta">
                        <span class="genre">{{ movie.genre }}</span>
                        <span class="year">{{ movie.year }}</span>
                    </div>
                    <div class="movie-rating">
                        <el-rate
                            v-model="movie.rating"
                            disabled
                            show-score
                            text-color="#ff9900"
                            score-template="{value}"
                        />
                    </div>
                    <p v-if="movie.comment" class="movie-comment">
                        {{ movie.comment }}
                    </p>
                </div>
                <div class="movie-actions">
                    <el-button
                        type="primary"
                        size="small"
                        @click="openEditDialog(movie)"
                        class="edit-btn"
                    >
                        <el-icon><Edit /></el-icon>
                        编辑
                    </el-button>
                    <el-button
                        type="danger"
                        size="small"
                        @click="deleteMovie(movie)"
                        class="delete-btn"
                    >
                        <el-icon><Delete /></el-icon>
                        删除
                    </el-button>
                </div>
            </div>

            <div v-if="movies.length === 0 && !loading" class="empty-state">
                <el-empty description="暂无影视，点击添加按钮开始记录" />
            </div>
        </div>

        <!-- Add/Edit Dialog -->
        <el-dialog
            v-model="dialogVisible"
            :title="isEditing ? '编辑影视' : '添加影视'"
            width="600px"
            :before-close="() => (dialogVisible = false)"
        >
            <el-form :model="form" label-width="80px">
                <el-form-item label="标题" required>
                    <el-input
                        v-model="form.title"
                        placeholder="请输入影视标题"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="导演" required>
                    <el-input
                        v-model="form.director"
                        placeholder="请输入导演"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="类型">
                    <el-input v-model="form.genre" placeholder="请输入影视类型" />
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
                    <el-input v-model="form.url" placeholder="请输入影视链接" />
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
                <el-button type="primary" @click="saveMovie">
                    {{ isEditing ? "更新" : "添加" }}
                </el-button>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.movies-page {
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

.movies-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
}

.movie-card {
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    overflow: hidden;
    background: white;
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
}

.movie-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.movie-image {
    height: 180px;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
}

.no-image {
    color: #ccc;
    font-size: 3rem;
}

.movie-info {
    padding: 16px;
    flex: 1;
}

.movie-title {
    margin: 0 0 8px 0;
    font-size: 1.2rem;
    font-weight: 600;
    color: #333;
    line-height: 1.4;
}

.movie-director {
    margin: 0 0 12px 0;
    color: #666;
    font-size: 0.9rem;
}

.movie-meta {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
    font-size: 0.85rem;
}

.genre {
    color: #e6a23c;
    background: #fdf6ec;
    padding: 2px 8px;
    border-radius: 4px;
}

.year {
    color: #909399;
}

.movie-rating {
    margin-bottom: 12px;
}

.movie-comment {
    margin: 8px 0;
    color: #666;
    font-size: 0.9rem;
    line-height: 1.4;
}

.movie-actions {
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

@media (max-width: 768px) {
    .movies-page {
        padding: 16px;
    }

    .page-header {
        flex-direction: column;
        gap: 16px;
        align-items: stretch;
    }

    .movies-grid {
        grid-template-columns: 1fr;
    }

    .movie-actions {
        flex-direction: column;
    }
}
</style>
