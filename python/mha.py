import torch
import torch.nn as nn
import torch.nn.functional as F


class MultiHeadAttention(nn.Module):
    def __init__(self, embed_dim, num_heads, dropout=0.1):
        super(MultiHeadAttention, self).__init__()

        assert embed_dim % num_heads == 0, "embed_dim must be divisible by num_heads"

        self.num_heads = num_heads
        self.head_dim = embed_dim // num_heads

        # Q, K, V 的线性映射
        self.q_linear = nn.Linear(embed_dim, embed_dim)
        self.k_linear = nn.Linear(embed_dim, embed_dim)
        self.v_linear = nn.Linear(embed_dim, embed_dim)

        # 最后的线性层
        self.out_linear = nn.Linear(embed_dim, embed_dim)

        # Dropout 层
        self.dropout = nn.Dropout(dropout)

        # Scale factor
        self.scale = self.head_dim ** 0.5

    def generate_causal_mask(self, seq_len):
        # 生成 causal mask (上三角矩阵)
        mask = torch.triu(torch.ones(seq_len, seq_len), diagonal=1)  # 上三角
        return mask.unsqueeze(0).unsqueeze(0)  # 增加 batch_size 和 num_heads 维度

    def generate_padding_mask(self, seq_len, padding_mask):
        # padding_mask: shape (batch_size, seq_len)
        # 生成 padding mask，确保填充部分不会影响注意力
        return padding_mask.unsqueeze(1).unsqueeze(2)  # (batch_size, 1, 1, seq_len)

    def forward(self, query, key, value, padding_mask=None):
        batch_size = query.size(0)
        seq_len = query.size(1)

        # 计算 Q, K, V
        Q = self.q_linear(query)  # (batch_size, seq_len, embed_dim)
        K = self.k_linear(key)  # (batch_size, seq_len, embed_dim)
        V = self.v_linear(value)  # (batch_size, seq_len, embed_dim)

        # 将 Q, K, V 拆分为 num_heads 个头
        Q = Q.view(batch_size, -1, self.num_heads, self.head_dim).transpose(1,
                                                                            2)  # (batch_size, num_heads, seq_len, head_dim)
        K = K.view(batch_size, -1, self.num_heads, self.head_dim).transpose(1,
                                                                            2)  # (batch_size, num_heads, seq_len, head_dim)
        V = V.view(batch_size, -1, self.num_heads, self.head_dim).transpose(1,
                                                                            2)  # (batch_size, num_heads, seq_len, head_dim)

        # 计算缩放点积注意力
        scores = torch.matmul(Q, K.transpose(-2, -1)) / self.scale  # (batch_size, num_heads, seq_len, seq_len)

        # 生成 Causal Mask (上三角 mask)
        causal_mask = self.generate_causal_mask(seq_len).to(scores.device)

        # 将 Causal Mask 应用到 scores
        scores = scores.masked_fill(causal_mask == 1, float('-inf'))

        # 如果有 Padding Mask，应用到 scores
        if padding_mask is not None:
            padding_mask = self.generate_padding_mask(seq_len, padding_mask).to(scores.device)
            scores = scores.masked_fill(padding_mask == 1, float('-inf'))

        # 计算注意力权重
        attn_weights = F.softmax(scores, dim=-1)  # (batch_size, num_heads, seq_len, seq_len)
        attn_weights = self.dropout(attn_weights)

        # 计算注意力输出
        attention_output = torch.matmul(attn_weights, V)  # (batch_size, num_heads, seq_len, head_dim)

        # 将头合并
        attention_output = attention_output.transpose(1, 2).contiguous().view(batch_size, -1,
                                                                              self.num_heads * self.head_dim)  # (batch_size, seq_len, embed_dim)

        # 通过输出线性层
        output = self.out_linear(attention_output)  # (batch_size, seq_len, embed_dim)

        return output, attn_weights


# 测试 MultiHeadAttention 层
batch_size = 2
seq_len = 4
embed_dim = 8
num_heads = 2

# 随机生成输入数据
query = torch.rand(batch_size, seq_len, embed_dim)
key = torch.rand(batch_size, seq_len, embed_dim)
value = torch.rand(batch_size, seq_len, embed_dim)

# 假设第一个 token 之后的 token 是 padding
padding_mask = torch.tensor([[0, 0, 1, 1], [0, 1, 0, 1]])  # (batch_size, seq_len)

# 创建 MultiHeadAttention 层
mha = MultiHeadAttention(embed_dim, num_heads)

# 前向传播
output, attn_weights = mha(query, key, value, padding_mask)

print("Output shape:", output.shape)  # 应该是 (batch_size, seq_len, embed_dim)
print("Attention Weights shape:", attn_weights.shape)  # 应该是 (batch_size, num_heads, seq_len, seq_len)