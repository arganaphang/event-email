import {
  Body,
  Button,
  Container,
  Head,
  Heading,
  Hr,
  Html,
  Link,
  Section,
  Tailwind,
  Text,
} from "@react-email/components";
import * as React from "react";

export const PaymentEmail = () => {
  return (
    <Html>
      <Head />
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans">
          <Container className="border border-solid border-[#eaeaea] rounded my-[40px] mx-auto p-[20px] w-[465px]">
            <Heading className="text-black text-[24px] font-normal text-center p-0 my-[30px] mx-0">
              Pay Invoice <strong>#{`{{ .ID }}`}</strong> for{" "}
              <strong>Our Services</strong>
            </Heading>
            <Text className="text-black text-[14px] leading-[24px]">
              Hello {`{{ .Name }}`},
            </Text>
            <Text className="text-black text-[14px] leading-[24px]">
              You has pending <strong>Invoice #{`{{ .ID }}`}</strong>. With
              total amount is{" "}
              <strong>
                <u>{`{{ .Amount }}`}</u>
              </strong>
            </Text>
            <Section className="text-center mt-[32px] mb-[32px]">
              <Button
                pX={20}
                pY={12}
                className="bg-[#000000] rounded text-white text-[12px] font-semibold no-underline text-center cursor-pointer"
              >
                Pay
              </Button>
            </Section>
            <Text className="text-black text-[14px] leading-[24px]">
              or copy and paste this URL into your browser:{" "}
              <Link
                href={`mailto:${`{{ .Email }}`}`}
                className="text-blue-600 no-underline"
              >
                Payment Page
              </Link>
            </Text>
            <Hr className="border border-solid border-[#eaeaea] my-[26px] mx-0 w-full" />
            <Text className="text-[#666666] text-[12px] leading-[24px]">
              This invoice was intended for{" "}
              <span className="text-black">{`{{ .Name }}`} </span>.This invoice
              was sent from{" "}
              <span className="text-black">sender@mailhog.local</span> located
              in <span className="text-black">Indonesia</span>. If you were not
              expecting this invoice, you can ignore this email. If you are
              concerned about your account's safety, please reply to this email
              to get in touch with us.
            </Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default PaymentEmail;
