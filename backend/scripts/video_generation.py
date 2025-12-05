# backend/scripts/video_generation.py
from modelscope.pipelines import pipeline
from modelscope.utils.constant import Tasks

pipe = pipeline(Tasks.text_to_video_synthesis, 
               model='damo/text-to-video-synthesis')

result = pipe({'text': 'A cat playing with a ball'})
result['output_video'].save('output.mp4')