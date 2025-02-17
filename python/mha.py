import torch
import torch.nn as nn
import torch.nn.functional as F
import math


class MultiHeadAttention(nn.Module):
    def __init__(self, embed_dim, num_heads, dropout=0.0, bias=True):
        """
        Args:
            embed_dim:  总嵌入维度 (必须能被 num_heads 整除)
            num_heads:  注意力头的数量
            dropout:    Attention权重 dropout概率
            bias:       是否在线性层使用偏置
        """
        super().__init__()
        assert embed_dim % num_heads == 0, "embed_dim must be divisible by num_heads"

        self.embed_dim = embed_dim
        self.num_heads = num_heads
        self.head_dim = embed_dim // num_heads  # 每个头的维度

        # 定义线性变换层 (Q, K, V 和输出投影)
        self.W_q = nn.Linear(embed_dim, embed_dim, bias=bias)
        self.W_k = nn.Linear(embed_dim, embed_dim, bias=bias)
        self.W_v = nn.Linear(embed_dim, embed_dim, bias=bias)
        self.out_proj = nn.Linear(embed_dim, embed_dim, bias=bias)

        self.dropout = nn.Dropout(dropout)

    def forward(self, q, k=None, v=None, key_padding_mask=None, attn_mask=None):
        """
        Args:
            q:                查询向量 (batch_size, seq_q, embed_dim)
            k:                键向量 (batch_size, seq_k, embed_dim) [可选，默认=q]
            v:                值向量 (batch_size, seq_v, embed_dim) [可选，默认=k]
            key_padding_mask: 键填充掩码 (batch_size, seq_k) [True表示填充位置]
            attn_mask:        注意力掩码 (seq_q, seq_k) 或 (batch_size, seq_q, seq_k)
                              [True表示需要屏蔽的位置 或 要加的值]

        Returns:
            output: 注意力输出 (batch_size, seq_q, embed_dim)
            attn_weights: 注意力权重 (batch_size, num_heads, seq_q, seq_k)
        """
        # 处理默认参数
        if k is None: k = q
        if v is None: v = k

        batch_size, seq_q, _ = q.size()
        seq_k, seq_v = k.size(1), v.size(1)

        # 步骤1: 线性变换 + 分头
        Q = self.W_q(q)  # (batch_size, seq_q, embed_dim)
        K = self.W_k(k)  # (batch_size, seq_k, embed_dim)
        V = self.W_v(v)  # (batch_size, seq_v, embed_dim)

        # 调整形状: (batch_size, num_heads, seq_len, head_dim)
        Q = Q.view(batch_size, seq_q, self.num_heads, self.head_dim).transpose(1, 2)
        K = K.view(batch_size, seq_k, self.num_heads, self.head_dim).transpose(1, 2)
        V = V.view(batch_size, seq_v, self.num_heads, self.head_dim).transpose(1, 2)

        # 步骤2: 计算注意力分数 (缩放点积)
        scores = torch.matmul(Q, K.transpose(-2, -1)) / math.sqrt(self.head_dim)

        # 步骤3: 应用注意力掩码
        if attn_mask is not None:
            if attn_mask.dtype == torch.bool:
                scores.masked_fill_(attn_mask, float('-inf'))  # 布尔掩码直接填充
            else:
                scores += attn_mask  # 加法掩码 (如因果掩码)

        # 步骤4: 应用键填充掩码
        if key_padding_mask is not None:
            # 扩展掩码维度: (batch_size, 1, 1, seq_k)
            mask = key_padding_mask.view(batch_size, 1, 1, seq_k)
            scores = scores.masked_fill(mask, float('-inf'))

        # 步骤5: 计算注意力权重 (softmax + dropout)
        attn_weights = F.softmax(scores, dim=-1)
        attn_weights = self.dropout(attn_weights)

        # 步骤6: 应用注意力权重到值向量
        output = torch.matmul(attn_weights, V)  # (batch_size, num_heads, seq_q, head_dim)

        # 步骤7: 合并多头输出
        output = output.transpose(1, 2).contiguous()  # (batch_size, seq_q, num_heads, head_dim)
        output = output.view(batch_size, seq_q, self.embed_dim)  # (batch_size, seq_q, embed_dim)

        # 步骤8: 输出投影
        output = self.out_proj(output)

        return output, attn_weights


# 使用示例
if __name__ == "__main__":
    embed_dim = 512
    num_heads = 8
    batch_size = 2
    seq_len = 10

    # 初始化模块
    mha = MultiHeadAttention(embed_dim, num_heads)

    # 创建测试输入
    q = torch.randn(batch_size, seq_len, embed_dim)

    # 带掩码的示例
    key_padding_mask = torch.zeros(batch_size, seq_len).bool()  # 假设第二个样本有填充
    key_padding_mask[1, 5:] = True

    # 前向传播
    output, attn_weights = mha(q, key_padding_mask=key_padding_mask)

    print("Output shape:", output.shape)  # 应输出 torch.Size([2, 10, 512])
    print("Attention weights shape:", attn_weights.shape)  # 应输出 torch.Size([2, 8, 10, 10])
