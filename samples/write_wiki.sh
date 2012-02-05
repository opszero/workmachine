curl -X PUT -H "Content-Type: application/json" \
    -d '{"program": "((:topic 
                        {:input [[:instructions \"topic_instruction\"]] 
                         :output [[:string \"topic\"]]}) 
                      (:research 
                        {:input [[:instructions \"research_instruction\"] 
                                 [:text \"topic\"]] 
                         :output [[:string \"info_link\"]]}) 
                      (:write 
                        {:input [[:instructions \"write_instruction\"] 
                                 [:text \"topic\"] 
                                 [:text \"info_link\"]] 
                         :output [[:text \"content\"]]}))",
         "jobs": [{"topic_instruction": "Give me a new topic."}]}' \
     http://localhost:3000/run

