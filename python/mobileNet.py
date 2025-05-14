import torch
import torch.nn as nn
from torchvision import models
from torch.utils.data import Dataset, DataLoader
from torchvision import transforms
from PIL import Image
import os


class AgeEstimationModel(nn.Module):
    def __init__(self):
        super(AgeEstimationModel, self).__init__()
        self.mobilenet = models.mobilenet_v2(pretrained=True)
        self.mobilenet.classifier = nn.Sequential(
            nn.Dropout(0.2),
            nn.Linear(self.mobilenet.last_channel, 1)  # 输出一个年龄值
        )

    def forward(self, x):
        return self.mobilenet(x)


transform = transforms.Compose([
    transforms.Resize((224, 224)),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                         std=[0.229, 0.224, 0.225])
])


def predict_age(model, image_path):
    model.eval()
    image = Image.open(image_path).convert('RGB')
    image = transform(image).unsqueeze(0).to("mps")
    with torch.no_grad():
        output = model(image)
    predicted_age = output.item()
    return predicted_age


# 示例
model = AgeEstimationModel().to("mps")
image_path = './test.jpeg'  # 替换为你的图片路径
predicted_age = predict_age(model, image_path)
print(f'Predicted Age: {predicted_age:.2f}')
