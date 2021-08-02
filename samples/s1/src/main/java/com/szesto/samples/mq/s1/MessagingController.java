package com.szesto.samples.mq.s1;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jms.core.JmsTemplate;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(path="/put", consumes = "text/plain")
@CrossOrigin(origins = "*")
public class MessagingController {

    private final JmsTemplate jmsTemplate;

    @Autowired
    public MessagingController(JmsTemplate jmsTemplate) {
        this.jmsTemplate = jmsTemplate;
    }

    @PostMapping(path="/message")
    public void putMessage(@RequestBody String message) {
        // inject from env
        final String qname = "Q.A";

        jmsTemplate.convertAndSend(qname, message);
    }
}
