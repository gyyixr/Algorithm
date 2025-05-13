# 导入相关需要的包
import math
import torch
import torch.nn as nn

import warnings

warnings.filterwarnings(action="ignore")

class SelfAttV4(nn.Module):
    def __init__(self, dim) -> None:
        super().__init__()
        self.dim = dim

        # 这样很清晰
        self.query_proj = nn.Linear(dim, dim)
        self.key_proj = nn.Linear(dim, dim)
        self.value_proj = nn.Linear(dim, dim)
        # 一般是 0.1 的 dropout，一般写作 config.attention_probs_dropout_prob
        # hidden_dropout_prob 一般也是 0.1
        self.att_drop = nn.Dropout(0.1)

        # 可以不写；具体和面试官沟通。
        # 这是 MultiHeadAttention 中的产物，这个留给 MultiHeadAttention 也没有问题；
        self.output_proj = nn.Linear(dim, dim)

    def forward(self, X, attention_mask=None):
        # attention_mask shape is: (batch, seq)
        # X shape is: (batch, seq, dim)

        Q = self.query_proj(X)
        K = self.key_proj(X)
        V = self.value_proj(X)

        att_weight = Q @ K.transpose(-1, -2) / math.sqrt(self.dim)
        if attention_mask is not None:
            # 给 weight 填充一个极小的值
            att_weight = att_weight.masked_fill(attention_mask == 0, float("-1e20"))

        att_weight = torch.softmax(att_weight, dim=-1)
        #print(att_weight)

        # 这里在 BERT中的官方代码也说很奇怪，但是原文中这么用了，所以继承了下来
        # （用于 output 后面会更符合直觉？）
        att_weight = self.att_drop(att_weight)

        output = att_weight @ V
        ret = self.output_proj(output)
        return ret


X = torch.rand(3, 4, 2)
b = torch.tensor(
    [
        [1, 1, 1, 0],
        [1, 1, 0, 0],
        [1, 0, 0, 0],
    ]
)
print(b.shape)
mask = b.unsqueeze(dim=1).repeat(1, 4, 1)

net = SelfAttV4(2)
print(net(X, mask).shape)