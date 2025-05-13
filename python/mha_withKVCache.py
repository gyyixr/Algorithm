import math
import torch
import torch.nn as nn


class MultiHeadAttentionWithKVCache(nn.Module):
    def __init__(self, hidden_dim, num_heads) -> None:
        super().__init__()
        self.num_heads = num_heads
        self.head_dim = hidden_dim // num_heads
        self.hidden_dim = hidden_dim

        self.q_proj = nn.Linear(hidden_dim, hidden_dim)
        self.k_proj = nn.Linear(hidden_dim, hidden_dim)
        self.v_proj = nn.Linear(hidden_dim, hidden_dim)

        self.att_dropout = nn.Dropout(0.1)
        self.o_proj = nn.Linear(hidden_dim, hidden_dim)

    def forward(self, X, attention_mask=None, kv_cache=None):
        batch_size, seq_len, _ = X.size()

        Q = self.q_proj(X)
        K = self.k_proj(X)
        V = self.v_proj(X)

        q = Q.view(batch_size, seq_len, self.num_heads, self.head_dim).transpose(1, 2)
        k = K.view(batch_size, seq_len, self.num_heads, self.head_dim).transpose(1, 2)
        v = V.view(batch_size, seq_len, self.num_heads, self.head_dim).transpose(1, 2)

        if kv_cache is not None:
            if 'k' in kv_cache and 'v' in kv_cache:
                k = torch.cat([kv_cache['k'], k], dim=2)
                v = torch.cat([kv_cache['v'], v], dim=2)
            kv_cache['k'], kv_cache['v'] = k, v

        attn_scores = torch.matmul(q, k.transpose(-2, -1)) / math.sqrt(self.head_dim)

        if attention_mask is not None:
            attn_scores = attn_scores.masked_fill(attention_mask == 0, float('-1e20'))

        attn_weights = torch.softmax(attn_scores, dim=-1)
        attn_weights = self.att_dropout(attn_weights)

        output = torch.matmul(attn_weights, v)
        output = output.transpose(1, 2).contiguous().view(batch_size, seq_len, -1)
        output = self.o_proj(output)
        if kv_cache is not None:
            return output, kv_cache
        else:
            return output



# 初始化模型和缓存
model = MultiHeadAttentionWithKVCache(hidden_dim=128, num_heads=8)
kv_cache = {}

# 模拟自回归生成过程
sequence = torch.rand(3, 10, 128)  # 假设输入序列长度为10
outputs = []

for t in range(sequence.size(1)):
    x_t = sequence[:, t:t + 1, :]  # 当前时间步的输入
    out, kv = model(x_t, kv_cache=kv_cache)
    print(kv['k'].size())
    print(kv['v'].size())
    outputs.append(out)

# 将所有时间步的输出拼接
final_output = torch.cat(outputs, dim=1)
print(final_output.size())
print(final_output)
print(kv_cache)

x = 5
y = 5
print(x == y)  # 输出: True，因为 x 和 y 的值相等

a = [1, 2, 3]
b = [1, 2, 3]
print(a == b)  # 输出: True，因为 a 和 b 的内容相同