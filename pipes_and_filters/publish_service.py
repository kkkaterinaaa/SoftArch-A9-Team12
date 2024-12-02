import json
import time
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
import os
from dotenv import load_dotenv

load_dotenv()


def send_email(subject, body, to_emails):
    from_email = os.getenv("EMAIL_MAIL")
    password = os.getenv("EMAIL_PASSWORD")

    msg = MIMEMultipart()
    msg['From'] = from_email
    msg['Subject'] = subject
    msg.attach(MIMEText(body, 'plain'))
    msg['To'] = ", ".join(to_emails)

    try:
        server = smtplib.SMTP('smtp.mail.ru', 587)
        server.starttls()
        server.login(from_email, password)
        text = msg.as_string()
        server.sendmail(from_email, to_emails, text)
        server.quit()
        print("[Publish Service] Email sent successfully!")
    except Exception as e:
        print(f"[Publish Service] Failed to send email: {e}")

def publish_message(output_queue):
    while True:
        message = output_queue.get()
        if message is None:
            break

        alias = message.get('alias')
        content = message.get('content')
        print(f"[Publish Service] Processing message from {alias}: {content}")
        to_emails = ["e.zaitseva@innopolis.university", "n.chekhonina@innopolis.university"]
        email_body = f"From: {alias}\nMessage: {content}"
        send_email("Pipes-and-filters", email_body, to_emails)
