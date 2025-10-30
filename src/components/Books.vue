<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { booksApi } from "../services/api";
import type { Book, DialogFormData } from "../services/types";

const books = ref<Book[]>([]);
const loading = ref(false);
const dialogVisible = ref(false);
const isEditing = ref(false);
const currentBook = ref<Book | null>(null);

const form = reactive<DialogFormData>({
    title: "",
    author: "",
    genre: "",
    year: new Date().getFullYear(),
    url: "",
    cover: "",
    comment: "",
    rating: 0,
});

// Load books
const loadBooks = async () => {
    loading.value = true;
    try {
        books.value = await booksApi.getBooks();
    } catch (error) {
        ElMessage.error("加载书籍失败");
        console.error(error);
    } finally {
        loading.value = false;
    }
};

// Open add dialog
const openAddDialog = () => {
    isEditing.value = false;
    currentBook.value = null;
    resetForm();
    dialogVisible.value = true;
};

// Open edit dialog
const openEditDialog = (book: Book) => {
    isEditing.value = true;
    currentBook.value = book;
    Object.assign(form, {
        title: book.title,
        author: book.author,
        genre: book.genre,
        year: book.year,
        url: book.url,
        cover: book.cover,
        comment: book.comment,
        rating: book.rating,
    });
    dialogVisible.value = true;
};

// Reset form
const resetForm = () => {
    Object.assign(form, {
        title: "",
        author: "",
        genre: "",
        year: new Date().getFullYear(),
        url: "",
        cover: "",
        comment: "",
        rating: 0,
    });
};

// Save book
const saveBook = async () => {
    if (!form.title || !form.author) {
        ElMessage.warning("请填写标题和作者");
        return;
    }

    try {
        const bookData: Book = {
            title: form.title,
            author: form.author!,
            genre: form.genre,
            year: form.year,
            url: form.url,
            cover: form.cover || "",
            comment: form.comment,
            rating: form.rating,
        };

        await booksApi.saveBook(bookData);
        ElMessage.success(isEditing.value ? "更新成功" : "添加成功");
        dialogVisible.value = false;
        await loadBooks();
    } catch (error) {
        ElMessage.error("保存失败");
        console.error(error);
    }
};

// Delete book
const deleteBook = async (book: Book) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除书籍 "${book.title}" 吗？`,
            "确认删除",
            {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning",
            },
        );

        await booksApi.deleteBook({
            title: book.title,
            author: book.author,
        });
        ElMessage.success("删除成功");
        await loadBooks();
    } catch (error) {
        if (error !== "cancel") {
            ElMessage.error("删除失败");
            console.error(error);
        }
    }
};

// Initialize
onMounted(() => {
    loadBooks();
});
</script>

<template>
    <div class="books-page">
        <div class="page-header">
            <h2>书籍收藏</h2>
            <el-button type="primary" @click="openAddDialog" class="add-button">
                <el-icon><Plus /></el-icon>
                添加书籍
            </el-button>
        </div>

        <div class="books-grid" v-loading="loading">
            <div
                v-for="book in books"
                :key="`${book.title}-${book.author}`"
                class="book-card"
            >
                <div class="book-image">
                    <img
                        v-if="book.cover"
                        :src="book.cover"
                        :alt="book.title"
                        @error="(e: any) => (e.target.style.display = 'none')"
                    />
                    <div v-else class="no-image">
                        <el-icon><Reading /></el-icon>
                    </div>
                </div>
                <div class="book-info">
                    <h3 class="book-title">{{ book.title }}</h3>
                    <p class="book-author">{{ book.author }}</p>
                    <div class="book-meta">
                        <span class="genre">{{ book.genre }}</span>
                        <span class="year">{{ book.year }}</span>
                    </div>
                    <div class="book-rating">
                        <el-rate
                            v-model="book.rating"
                            disabled
                            show-score
                            text-color="#ff9900"
                            score-template="{value}"
                        />
                    </div>
                    <p v-if="book.comment" class="book-comment">
                        {{ book.comment }}
                    </p>
                </div>
                <div class="book-actions">
                    <el-button
                        type="primary"
                        size="small"
                        @click="openEditDialog(book)"
                        class="edit-btn"
                    >
                        <el-icon><Edit /></el-icon>
                        编辑
                    </el-button>
                    <el-button
                        type="danger"
                        size="small"
                        @click="deleteBook(book)"
                        class="delete-btn"
                    >
                        <el-icon><Delete /></el-icon>
                        删除
                    </el-button>
                </div>
            </div>

            <div v-if="books===null || (books.length === 0 && !loading)" class="empty-state">
                <el-empty description="暂无书籍，点击添加按钮开始记录" />
            </div>
        </div>

        <!-- Add/Edit Dialog -->
        <el-dialog
            v-model="dialogVisible"
            :title="isEditing ? '编辑书籍' : '添加书籍'"
            width="600px"
            :before-close="() => (dialogVisible = false)"
        >
            <el-form :model="form" label-width="80px">
                <el-form-item label="标题" required>
                    <el-input
                        v-model="form.title"
                        placeholder="请输入书籍标题"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="作者" required>
                    <el-input
                        v-model="form.author"
                        placeholder="请输入作者"
                        :disabled="isEditing"
                    />
                </el-form-item>
                <el-form-item label="类型">
                    <el-input v-model="form.genre" placeholder="请输入书籍类型" />
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
                    <el-input v-model="form.url" placeholder="请输入书籍链接" />
                </el-form-item>
                <el-form-item label="封面">
                    <el-input
                        v-model="form.cover"
                        placeholder="请输入封面图片链接"
                    />
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
                <el-button type="primary" @click="saveBook">
                    {{ isEditing ? "更新" : "添加" }}
                </el-button>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.books-page {
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

.books-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
}

.book-card {
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    overflow: hidden;
    background: white;
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
}

.book-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.book-image {
    height: 180px;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
}

.book-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.no-image {
    color: #ccc;
    font-size: 3rem;
}

.book-info {
    padding: 16px;
    flex: 1;
}

.book-title {
    margin: 0 0 8px 0;
    font-size: 1.2rem;
    font-weight: 600;
    color: #333;
    line-height: 1.4;
}

.book-author {
    margin: 0 0 12px 0;
    color: #666;
    font-size: 0.9rem;
}

.book-meta {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
    font-size: 0.85rem;
}

.genre {
    color: #67c23a;
    background: #f0f9eb;
    padding: 2px 8px;
    border-radius: 4px;
}

.year {
    color: #909399;
}

.book-rating {
    margin-bottom: 12px;
}

.book-comment {
    margin: 8px 0;
    color: #666;
    font-size: 0.9rem;
    line-height: 1.4;
}

.book-actions {
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
    .books-page {
        padding: 16px;
    }

    .page-header {
        flex-direction: column;
        gap: 16px;
        align-items: stretch;
    }

    .books-grid {
        grid-template-columns: 1fr;
    }

    .book-actions {
        flex-direction: column;
    }
}
</style>
